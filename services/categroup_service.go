package services

import (
	"github.com/xiaoxuan6/gsv/pkg/global"
	"slices"
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

	languagesStarRepos := make(map[string][]*global.GRepository, 0)
	for _, repo := range repos {
		repo := repo
		wg.Add(1)
		go func() {
			defer wg.Done()

			language := repo.Repository.GetLanguage()
			if gRepository, ok := languagesStarRepos[language]; ok {
				gRepository = append(gRepository, repo)
				languagesStarRepos[language] = gRepository
			} else {
				languagesStarRepos[language] = []*global.GRepository{repo}
			}
		}()
	}

	wg.Wait()
	global.AccountsLanguageStarRepose[global.CurrentAccount] = languagesStarRepos
}
