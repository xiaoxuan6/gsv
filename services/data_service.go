package services

import (
	"fmt"
	"github-stars/pkg/github"
	"github-stars/pkg/global"
	"github-stars/pkg/translate"
	github2 "github.com/google/go-github/v64/github"
	"strings"
	"sync"
)

var (
	wg   sync.WaitGroup
	lock sync.Mutex
)

func FetchDataWithPage(account string, page int) (items []string, nextPage int) {
	return FetchData(account, page, 100)
}

func FetchData(account string, page, perPage int) (items []string, nextPage int) {
	allStarRepos, nextPage := github.AllStarsRepos(account, page, perPage)

	var allRepos []*github2.Repository
	for _, starRepos := range allStarRepos {
		repos := starRepos.Repository
		allRepos = append(allRepos, repos)
	}

	items = CheckRepos(allRepos)

	return
}

func CheckRepos(repos []*github2.Repository) (items []string) {
	for _, val := range repos {
		val := val
		wg.Add(1)
		go func() {
			defer wg.Done()

			lock.Lock()
			desc := checkItem(val)
			NewRepository := &global.GRepository{
				Repository:    val,
				DescriptionZh: desc,
			}

			var AllStarRepos []*global.GRepository
			if value, ok := global.AccountsAllStarRepos[global.CurrentAccount]; ok {
				AllStarRepos = append(value, NewRepository)
			} else {
				AllStarRepos = append(AllStarRepos, NewRepository)
			}
			global.AccountsAllStarRepos[global.CurrentAccount] = AllStarRepos

			items = append(items, desc)
			lock.Unlock()
		}()
	}
	wg.Wait()

	return
}

func checkItem(repos *github2.Repository) string {
	var fullname string
	if repos != nil && repos.FullName != nil {
		fullname = repos.GetFullName()
	}

	var language string
	if repos != nil && repos.Language != nil {
		language = repos.GetLanguage()
	}

	var description string
	if repos != nil && repos.Description != nil {
		description = translate.Translation(repos.GetDescription())
	}
	description = strings.ReplaceAll(strings.ReplaceAll(description, " | ", ""), "|", "")

	return fmt.Sprintf("【%s】（%s） - %s", fullname, description, language)
}

func CheckItem(repos []*global.GRepository) (items []string) {
	for _, val := range repos {
		lock.Lock()
		items = append(items, val.DescriptionZh)
		lock.Unlock()
	}

	return
}