package github

import (
	"context"
	"fmt"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v64/github"
	"os"
)

var (
	client *github.Client

	ctx = context.Background()
)

func newGithubClient() *github.Client {
	if client != nil {
		return client
	}

	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		panic(err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if os.Getenv("GITHUB_TOKEN") == "" {
		panic("github token empty")
	}

	client = github.NewClient(rateLimiter).WithAuthToken(token)
	return client
}

func AllStarsRepos(user string, page, perPage int) (allStarsRepos []*github.StarredRepository, nextPage int) {
	opt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	repos, resp, err := newGithubClient().Activity.ListStarred(ctx, user, opt)
	if err != nil {
		fmt.Println("fetch git stars error: ", err.Error())
		return
	}

	allStarsRepos = append(allStarsRepos, repos...)
	if resp.NextPage != 0 {
		nextPage = resp.NextPage
	}

	return
}

func SearchOwner(username string) []*github.User {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 10,
		},
	}
	result, _, err := newGithubClient().Search.Users(ctx, username, opt)
	if err != nil {
		fmt.Println("fetch github owner error: ", err.Error())
		return nil
	}

	return result.Users
}
