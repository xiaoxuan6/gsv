package services

import (
	"fmt"
	github2 "github.com/google/go-github/v64/github"
	"github.com/xiaoxuan6/gsv/pkg/github"
	"github.com/xiaoxuan6/gsv/pkg/global"
	"github.com/xiaoxuan6/gsv/pkg/translate"
	"strings"
	"sync"
)

var (
	wg                   sync.WaitGroup
	lock                 sync.Mutex
	accountsAllStarRepos sync.Map
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

	if len(allRepos) > 0 {
		items = CheckRepos(allRepos)
	} else {
		allStarsRepos := make([]*global.GRepository, 0)
		global.AccountsAllStarRepos[global.CurrentAccount] = allStarsRepos
	}
	global.AccountsStarReposNextPage[global.CurrentAccount] = nextPage

	return
}

func CheckRepos(repos []*github2.Repository) (items []string) {
	for _, val := range repos {
		val := val
		wg.Add(1)
		go func() {
			defer wg.Done()

			desc := checkItem(val)
			newRepository := &global.GRepository{
				Repository:    val,
				DescriptionZh: desc,
			}

			if v, ok := accountsAllStarRepos.Load(global.CurrentAccount); ok {
				allStarRepos := v.([]*global.GRepository)
				allStarRepos = append(allStarRepos, newRepository)
				accountsAllStarRepos.Store(global.CurrentAccount, allStarRepos)
			} else {
				accountsAllStarRepos.Store(global.CurrentAccount, []*global.GRepository{newRepository})
			}

			items = append(items, desc)
		}()
	}
	wg.Wait()

	if v, ok := accountsAllStarRepos.Load(global.CurrentAccount); ok {
		allStarRepos := v.([]*global.GRepository)
		global.AccountsAllStarRepos[global.CurrentAccount] = allStarRepos
	}

	return
}

func checkItem(repos *github2.Repository) string {
	var fullname, language, description string
	if repos != nil {
		if repos.FullName != nil {
			fullname = repos.GetFullName()
		}
		if repos.Language != nil {
			language = repos.GetLanguage()
		}
		if repos.Description != nil {
			description = translate.Translation(repos.GetDescription())
		}
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
