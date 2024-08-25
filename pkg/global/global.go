package global

import github2 "github.com/google/go-github/v64/github"

type GRepository struct {
	Repository    *github2.Repository
	DescriptionZh string `json:"description_zh"`
}

var (
	CurrentAccount string

	AccountsAllStarRepos = make(map[string][]*GRepository, 0)
)