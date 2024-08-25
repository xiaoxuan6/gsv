package github

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v64/github"
	"github.com/mitchellh/go-homedir"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

var (
	client *github.Client

	ctx = context.Background()
)

func fetchToken() string {
	gitCredentials, err := homedir.Expand("~/.git-credentials")
	if err != nil {
		return ""
	}

	body, err := ioutil.ReadFile(gitCredentials)
	if err != nil {
		return ""
	}

	r := bufio.NewReader(strings.NewReader(string(body)))
	for {
		line, _, err1 := r.ReadLine()
		if err1 == io.EOF {
			break
		}

		if strings.HasSuffix(string(line), "@github.com") {
			u, _ := url.Parse(string(line))
			if password, ok := u.User.Password(); ok && strings.HasPrefix(password, "ghp_") {
				return password
			}
		}
	}

	return ""
}

func newGithubClient() *github.Client {
	if client != nil {
		return client
	}

	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		fmt.Printf("github ratelimt error：%s", err.Error())
		return nil
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = fetchToken()
		if token == "" {
			fmt.Printf(`
github token empty.
Please set: export GITHUB_TOKEN="xxxx" or create .env file
github token：https://github.com/settings/tokens

`)
			return nil
		}
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
