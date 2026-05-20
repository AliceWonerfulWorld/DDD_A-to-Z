package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	mypageapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/mypage"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
)

type githubUser struct {
	PublicRepos int    `json:"public_repos"`
	CreatedAt   string `json:"created_at"`
}

type repoItem struct {
	StargazersCount int  `json:"stargazers_count"`
	Fork            bool `json:"fork"`
}

type searchIssuesResponse struct {
	TotalCount int `json:"total_count"`
}

func (c *RepositoryClient) FetchStats(ctx context.Context, accessToken, username string) (*mypageapp.GitHubStats, error) {
	var (
		user        *githubUser
		repos       []repoItem
		totalPRs    int
		totalIssues int
		yearly      *yearlyStats
	)

	errs := make(chan error, 5)

	go func() {
		var err error
		user, err = c.fetchUser(ctx, accessToken, username)
		errs <- err
	}()
	go func() {
		var err error
		repos, err = c.fetchAllRepos(ctx, accessToken, username)
		errs <- err
	}()
	go func() {
		var err error
		totalPRs, err = c.searchIssueCount(ctx, accessToken, fmt.Sprintf("author:%s+type:pr", username))
		errs <- err
	}()
	go func() {
		var err error
		totalIssues, err = c.searchIssueCount(ctx, accessToken, fmt.Sprintf("author:%s+type:issue", username))
		errs <- err
	}()
	go func() {
		var err error
		yearly, err = c.fetchYearlyContributions(ctx, accessToken, username)
		errs <- err
	}()

	for i := 0; i < 5; i++ {
		if e := <-errs; e != nil {
			return nil, e
		}
	}

	totalStars := 0
	nonForkCount := 0
	for _, r := range repos {
		totalStars += r.StargazersCount
		if !r.Fork {
			nonForkCount++
		}
	}

	yearlyCommits := 0
	yearlyContributions := 0
	if yearly != nil {
		yearlyCommits = yearly.totalCommits
		yearlyContributions = yearly.contributionDays
	}

	return &mypageapp.GitHubStats{
		TotalStars:          totalStars,
		TotalPRs:            totalPRs,
		TotalIssues:         totalIssues,
		ContributedTo:       nonForkCount,
		PublicRepos:         user.PublicRepos,
		GitHubCreatedAt:     user.CreatedAt,
		YearlyCommits:       yearlyCommits,
		YearlyContributions: yearlyContributions,
	}, nil
}

func (c *RepositoryClient) fetchUser(ctx context.Context, accessToken, username string) (*githubUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/users/"+username, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var user githubUser
	if err := decodeResponse(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *RepositoryClient) fetchAllRepos(ctx context.Context, accessToken, username string) ([]repoItem, error) {
	var all []repoItem
	nextURL := c.baseURL + "/users/" + username + "/repos?per_page=100&type=owner"

	for nextURL != "" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, nextURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("User-Agent", userAgent)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		var page []repoItem
		if err := decodeResponse(resp, &page); err != nil {
			return nil, err
		}
		all = append(all, page...)
		nextURL = nextLink(resp.Header.Get("Link"))
	}
	return all, nil
}

func (c *RepositoryClient) searchIssueCount(ctx context.Context, accessToken, query string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/search/issues?q="+query+"&per_page=1", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}

	var payload searchIssuesResponse
	if err := decodeResponse(resp, &payload); err != nil {
		return 0, err
	}
	return payload.TotalCount, nil
}

type yearlyStats struct {
	totalCommits     int
	contributionDays int
	daily            []repositoryanalysis.DailyContribution
}

func (c *RepositoryClient) fetchYearlyContributions(ctx context.Context, accessToken, username string) (*yearlyStats, error) {
	year := time.Now().UTC().Year()
	query := fmt.Sprintf(`{
		user(login: "%s") {
			contributionsCollection(from: "%d-01-01T00:00:00Z", to: "%d-12-31T23:59:59Z") {
				totalCommitContributions
				contributionCalendar {
					totalContributions
					weeks {
						contributionDays {
							date
							contributionCount
						}
					}
				}
			}
		}
	}`, username, year, year)

	body := map[string]string{
		"query": query,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.github.com/graphql", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			User struct {
				ContributionsCollection struct {
					TotalCommitContributions int `json:"totalCommitContributions"`
					ContributionCalendar     struct {
						TotalContributions int `json:"totalContributions"`
						Weeks              []struct {
							ContributionDays []struct {
								Date  string `json:"date"`
								Count int    `json:"contributionCount"`
							} `json:"contributionDays"`
						} `json:"weeks"`
					} `json:"contributionCalendar"`
				} `json:"contributionsCollection"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	coll := result.Data.User.ContributionsCollection

	dailyMap := map[string]int{}

	for _, week := range coll.ContributionCalendar.Weeks {
		for _, day := range week.ContributionDays {
			dailyMap[day.Date] = day.Count
		}
	}

	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	var daily []repositoryanalysis.DailyContribution
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		daily = append(daily, repositoryanalysis.DailyContribution{
			Date:  d,
			Count: dailyMap[dateKey],
		})
	}

	return &yearlyStats{
		totalCommits:     coll.TotalCommitContributions,
		contributionDays: coll.ContributionCalendar.TotalContributions,
		daily:            daily,
	}, nil
}
