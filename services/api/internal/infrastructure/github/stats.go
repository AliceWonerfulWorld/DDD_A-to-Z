package github

import (
	"context"
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
		daily       []repositoryanalysis.DailyContribution
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
		since := time.Now().AddDate(0, 0, -364)
		daily, err = c.listUserContributions(ctx, accessToken, username, since)
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
	for _, d := range daily {
		yearlyCommits += d.Count
	}

	return &mypageapp.GitHubStats{
		TotalStars:          totalStars,
		TotalPRs:            totalPRs,
		TotalIssues:         totalIssues,
		ContributedTo:       nonForkCount,
		PublicRepos:         user.PublicRepos,
		GitHubCreatedAt:     user.CreatedAt,
		YearlyCommits:       yearlyCommits,
		YearlyContributions: len(daily),
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

func (c *RepositoryClient) listUserContributions(ctx context.Context, accessToken, username string, since time.Time) ([]repositoryanalysis.DailyContribution, error) {
	page := 1
	perPage := 100
	dateCounts := map[string]int{}

	for {
		requestURL := fmt.Sprintf("%s/search/commits?q=author:%s+committer-date:>=%s&per_page=%d&page=%d",
			c.baseURL, username, since.Format("2006-01-02"), perPage, page)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("X-GitHub-Api-Version", gitHubAPIVersion)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		var payload struct {
			TotalCount int `json:"total_count"`
			Items      []struct {
				Commit struct {
					Author struct {
						Date time.Time `json:"date"`
					} `json:"author"`
				} `json:"commit"`
			} `json:"items"`
		}
		if err := decodeResponse(resp, &payload); err != nil {
			return nil, err
		}

		for _, item := range payload.Items {
			date := item.Commit.Author.Date.Format("2006-01-02")
			dateCounts[date]++
		}

		if len(payload.Items) < perPage || page >= 10 {
			break
		}
		page++
	}

	var contributions []repositoryanalysis.DailyContribution
	start := since
	end := time.Now()
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		contributions = append(contributions, repositoryanalysis.DailyContribution{
			Date:  d,
			Count: dateCounts[dateKey],
		})
	}

	return contributions, nil
}
