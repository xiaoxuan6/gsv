package services

import (
	"github.com/xiaoxuan6/gsv/pkg/global"
	"slices"
	"sync"
)

func Category(repos []*global.GRepository) {
	languages := make([]string, 0)
	for _, repo := range repos {
		language := repo.Repository.GetLanguage()
		if ok := slices.Contains(languages, language); !ok {
			languages = append(languages, language)
		}
	}
	global.AccountsAllLanguages[global.CurrentAccount] = languages

	var languagesMap sync.Map
	for _, repo := range repos {
		repo := repo
		wg.Add(1)
		go func() {
			defer wg.Done()

			language := repo.Repository.GetLanguage()
			if v, ok := languagesMap.Load(language); ok {
				gRepository := v.([]*global.GRepository)
				gRepository = append(gRepository, repo)
				languagesMap.Store(language, gRepository)
			} else {
				languagesMap.Store(language, []*global.GRepository{repo})
			}
		}()
	}
	wg.Wait()

	languagesStarRepos := make(map[string][]*global.GRepository, 0)
	for _, language := range languages {
		value, _ := languagesMap.Load(language)
		languagesStarRepos[language] = value.([]*global.GRepository)
	}
	global.AccountsLanguageStarRepose[global.CurrentAccount] = languagesStarRepos
}
