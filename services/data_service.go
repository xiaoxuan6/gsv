package services

import (
	"fmt"
	github2 "github.com/google/go-github/v64/github"
	"github.com/xiaoxuan6/gsv/pkg/github"
	"github.com/xiaoxuan6/gsv/pkg/global"
	"github.com/xiaoxuan6/gsv/pkg/translate"
	"slices"
	"strings"
	"sync"
)

var (
	wg   sync.WaitGroup
	lock sync.Mutex
)

func FetchDataWithPage(account string, page int) (items []string, nextPage int) {
	return FetchData(account, page, global.PageCount)
}

func FetchData(account string, page, perPage int) (items []string, nextPage int) {
	allRepos, nextPage := FetchStarRepos(account, page, perPage)
	if len(allRepos) > 0 {
		items = CheckRepos(allRepos)
	} else {
		allStarsRepos := make([]*global.GRepository, 0)
		global.AccountsAllStarRepos[global.CurrentAccount] = allStarsRepos
	}
	global.AccountsStarReposNextPage[global.CurrentAccount] = nextPage

	return
}

func FetchStarRepos(account string, page, perPage int) ([]*github2.Repository, int) {
	allStarRepos, nextPage := github.AllStarsRepos(account, page, perPage)

	var allRepos []*github2.Repository
	for _, starRepos := range allStarRepos {
		repos := starRepos.Repository
		if repos.GetDescription() == "" || repos.GetLanguage() == "" {
			continue
		}

		allRepos = append(allRepos, repos)
	}

	return allRepos, nextPage
}

func CheckRepos(repos []*github2.Repository) (items []string) {
	var accountsAllStarRepos sync.Map
	for _, val := range repos {
		val := val
		wg.Add(1)
		go func() {
			defer wg.Done()

			if slices.Contains(History, val.GetFullName()) {
				fmt.Println(fmt.Sprintf("repo %s exists", val.GetFullName()))
				return
			}

			desc, descT, stat := checkItem(val)
			newRepository := &global.GRepository{
				Repository:           val,
				DescriptionZh:        desc,
				DescriptionTranslate: descT,
				TranslateStat:        stat,
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

	go Category(global.AccountsAllStarRepos[global.CurrentAccount])

	return
}

func checkItem(repos *github2.Repository) (string, string, bool) {
	var fullname, language, description string
	var ok bool
	if repos != nil {
		if fullName := repos.GetFullName(); fullName != "" {
			fullname = fullName
		}
		if gLanguage := repos.GetLanguage(); gLanguage != "" {
			language = gLanguage
		}
		if gDescription := repos.GetDescription(); gDescription != "" {
			description, ok = translate.Translation(gDescription)
		}
	}

	description = strings.ReplaceAll(strings.ReplaceAll(description, " | ", ""), "|", "")
	return fmt.Sprintf("【%s】（%s） - %s", fullname, description, language), description, ok
}

func CheckItem(repos []*global.GRepository) (items []string) {
	for _, val := range repos {
		lock.Lock()
		items = append(items, val.DescriptionZh)
		lock.Unlock()
	}

	return
}
