package repositoryanalysis

import "time"

type CommitItem struct {
	SHA       string
	Message   string
	Committed time.Time
}

type PullRequestItem struct {
	Number    int
	Title     string
	CreatedAt time.Time
}

type LanguageContribution struct {
	Name string
	CP   int64
	SP   int64
}

type Contribution struct {
	Repo      string
	Type      string // "commit" or "pull_request"
	Message   string
	Language  string
	CP        int64
	Timestamp time.Time
}

type AnalysisResult struct {
	TotalCommits      int64
	TotalPRs          int64
	TotalCP           int64
	LanguageBreakdown []LanguageContribution
	Contributions     []Contribution
}
