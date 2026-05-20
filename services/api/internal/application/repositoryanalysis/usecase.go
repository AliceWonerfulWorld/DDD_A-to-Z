package repositoryanalysis

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

const analysisPeriod = -30 * 24 * time.Hour

var (
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrMissingGitHubToken = errors.New("github token is missing")
)

const prCP = 5

type UseCase struct {
	current   CurrentUserRepository
	tokens    TokenRepository
	repos     RepositoryClient
	repoStore RepositoryStore
	logs      ContributionStore
	commits   GitHubCommitClient
	prs       GitHubPRClient
	languages GitHubLanguageClient
	cp        CPEarner
	sp        SPEarner
	cpBalance CPBalanceProvider
	now       func() time.Time
}

func NewUseCase(
	current CurrentUserRepository,
	tokens TokenRepository,
	repos RepositoryClient,
	repoStore RepositoryStore,
	commits GitHubCommitClient,
	prs GitHubPRClient,
	languages GitHubLanguageClient,
	cp CPEarner,
	sp SPEarner,
	cpBalance CPBalanceProvider,
) *UseCase {
	return NewUseCaseWithContributionStore(current, tokens, repos, repoStore, nil, commits, prs, languages, cp, sp, cpBalance)
}

func NewUseCaseWithContributionStore(
	current CurrentUserRepository,
	tokens TokenRepository,
	repos RepositoryClient,
	repoStore RepositoryStore,
	logs ContributionStore,
	commits GitHubCommitClient,
	prs GitHubPRClient,
	languages GitHubLanguageClient,
	cp CPEarner,
	sp SPEarner,
	cpBalance CPBalanceProvider,
) *UseCase {
	return &UseCase{
		current:   current,
		tokens:    tokens,
		repos:     repos,
		repoStore: repoStore,
		logs:      logs,
		commits:   commits,
		prs:       prs,
		languages: languages,
		cp:        cp,
		sp:        sp,
		cpBalance: cpBalance,
		now:       time.Now,
	}
}

type AnalysisResult struct {
	TotalCommits      int64                                     `json:"totalCommits"`
	TotalPRs          int64                                     `json:"totalPRs"`
	TotalCP           int64                                     `json:"totalCP"`
	TotalBalance      int64                                     `json:"totalBalance"`
	LanguageBreakdown []repositoryanalysis.LanguageContribution `json:"languageBreakdown"`
	Contributions     []repositoryanalysis.Contribution         `json:"contributions"`
}

func (u *UseCase) Analyze(ctx context.Context, sessionToken string) (AnalysisResult, error) {
	now := u.now()

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, now)
	if err != nil {
		return AnalysisResult{}, err
	}
	if !ok {
		return AnalysisResult{}, ErrUnauthenticated
	}

	return u.AnalyzeForUser(ctx, appUser, sessionToken, now)
}

func (u *UseCase) AnalyzeForUser(ctx context.Context, appUser user.User, sessionToken string, now time.Time) (AnalysisResult, error) {
	accessToken, ok, err := u.tokens.GitHubAccessToken(ctx, appUser.ID)
	if err != nil {
		return AnalysisResult{}, err
	}
	if !ok {
		return AnalysisResult{}, ErrMissingGitHubToken
	}

	syncedAt := now
	repos, err := u.repos.ListRepositories(ctx, accessToken, appUser.ID, syncedAt)
	if err != nil {
		return AnalysisResult{}, err
	}

	if err := u.repoStore.UpsertRepositories(ctx, repos); err != nil {
		return AnalysisResult{}, err
	}

	since := now.Add(analysisPeriod)
	lastAnalyzedAt, err := u.cpBalance.GetLastAnalyzedAt(ctx, appUser.ID)
	if err != nil {
		return AnalysisResult{}, err
	}
	if lastAnalyzedAt != nil && lastAnalyzedAt.After(since) {
		since = *lastAnalyzedAt
	}

	var totalCommits int64
	var totalPRs int64
	langCP := map[string]int64{}
	var contributions []repositoryanalysis.Contribution
	username := appUser.GitHubAccount.Username
	var apiErr bool

	for _, repo := range repos {
		if repo.Fork || repo.Archived {
			continue
		}
		if repo.GitHubUpdatedAt.Before(since) {
			continue
		}

		commits, err := u.commits.ListCommits(ctx, accessToken, repo.Owner, repo.Name, username, since)
		if err != nil {
			commits = nil
			apiErr = true
		}

		prs, err := u.prs.ListPullRequests(ctx, accessToken, repo.Owner, repo.Name, username, since)
		if err != nil {
			prs = nil
			apiErr = true
		}

		if len(commits) == 0 && len(prs) == 0 {
			continue
		}

		langs, err := u.languages.ListLanguages(ctx, accessToken, repo.Owner, repo.Name)
		if err != nil {
			langs = map[string]int64{repo.Language: 1}
			apiErr = true
		} else if len(langs) == 0 {
			langs = map[string]int64{repo.Language: 1}
		}

		totalBytes := int64(0)
		for _, bytes := range langs {
			totalBytes += bytes
		}
		if totalBytes == 0 {
			langs = map[string]int64{repo.Language: 1}
			totalBytes = 1
		}

		repoCP := int64(len(commits)) + int64(len(prs))*prCP
		totalCommits += int64(len(commits))
		totalPRs += int64(len(prs))

		type langAlloc struct {
			name string
			raw  float64
		}
		var allocs []langAlloc
		var allocated int64
		for lang, bytes := range langs {
			raw := float64(repoCP) * float64(bytes) / float64(totalBytes)
			cp := int64(raw)
			if cp > 0 {
				langCP[lang] += cp
				allocated += cp
			}
			allocs = append(allocs, langAlloc{name: lang, raw: raw - float64(cp)})
		}
		remaining := repoCP - allocated
		for remaining > 0 {
			best := 0
			for i := 1; i < len(allocs); i++ {
				if allocs[i].raw > allocs[best].raw {
					best = i
				}
			}
			if allocs[best].raw <= 0 {
				break
			}
			langCP[allocs[best].name]++
			allocs[best].raw = -1
			remaining--
		}

		primaryLang := repo.Language
		if primaryLang == "" {
			for lang := range langs {
				primaryLang = lang
				break
			}
		}

		for _, c := range commits {
			contributions = append(contributions, repositoryanalysis.Contribution{
				Repo:       repo.FullName,
				Type:       "commit",
				ExternalID: c.SHA,
				Message:    c.Message,
				Language:   primaryLang,
				CP:         1,
				Timestamp:  c.Committed,
			})
		}
		for _, pr := range prs {
			contributions = append(contributions, repositoryanalysis.Contribution{
				Repo:       repo.FullName,
				Type:       "pull_request",
				ExternalID: prExternalID(pr.Number),
				Message:    pr.Title,
				Language:   primaryLang,
				CP:         prCP,
				Timestamp:  pr.CreatedAt,
			})
		}
	}

	totalCP := int64(0)
	for _, cp := range langCP {
		totalCP += cp
	}

	langSP := make(map[string]int64, len(langCP))

	if !apiErr {
		if u.logs != nil && len(contributions) > 0 {
			if err := u.logs.UpsertAnalysisContributions(ctx, appUser.ID, contributions, now); err != nil {
				return AnalysisResult{}, err
			}
		}
		if totalCP > 0 {
			if err := u.cp.Earn(ctx, appUser.ID, totalCP, "contribution analysis reward", "analysis", "initial"); err != nil {
				return AnalysisResult{}, err
			}
		}
		for lang, cp := range langCP {
			if err := u.sp.EarnSP(ctx, appUser.ID, lang, cp, "skill point from contribution analysis", "analysis", "initial"); err != nil {
				if errors.Is(err, contributionpointapp.ErrUnsupportedPointType) {
					continue
				}
				return AnalysisResult{}, err
			}
			langSP[lang] = cp
		}
		if err := u.cpBalance.UpdateLastAnalyzedAt(ctx, appUser.ID, now); err != nil {
			return AnalysisResult{}, err
		}
	}

	breakdown := make([]repositoryanalysis.LanguageContribution, 0, len(langCP))
	for name, cp := range langCP {
		breakdown = append(breakdown, repositoryanalysis.LanguageContribution{Name: name, CP: cp, SP: langSP[name]})
	}
	sort.Slice(breakdown, func(i, j int) bool {
		return breakdown[i].CP > breakdown[j].CP
	})

	balance, err := u.cpBalance.GetBalance(ctx, appUser.ID)
	if err != nil {
		return AnalysisResult{}, err
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].Timestamp.After(contributions[j].Timestamp)
	})
	if len(contributions) > 50 {
		contributions = contributions[:50]
	}

	return AnalysisResult{
		TotalCommits:      totalCommits,
		TotalPRs:          totalPRs,
		TotalCP:           totalCP,
		TotalBalance:      balance,
		LanguageBreakdown: breakdown,
		Contributions:     contributions,
	}, nil
}

func prExternalID(number int) string {
	return "PR#" + strconv.Itoa(number)
}
