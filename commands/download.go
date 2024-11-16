package commands

import (
	"bufio"
	"errors"
	"fmt"
	github2 "github.com/google/go-github/v64/github"
	"github.com/urfave/cli/v2"
	"github.com/xiaoxuan6/gsv/pkg/github"
	"github.com/xiaoxuan6/gsv/pkg/global"
	"github.com/xiaoxuan6/gsv/pkg/translate"
	"github.com/xiaoxuan6/gsv/services"
	"github.com/xuri/excelize/v2"
	"io"
	"mvdan.cc/xurls/v2"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

func Download() *cli.Command {
	return &cli.Command{
		Name:        "download",
		Usage:       "下载指定用户的所有 start repos",
		Description: "下载指定用户的所有 start repos",
		Aliases:     []string{"d"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "account",
				Aliases:  []string{"a"},
				Required: true,
				Usage:    "需要查找的 github 账号，支持多个用户：a|b|c",
			},
			&cli.BoolFlag{
				Name:    "filter",
				Aliases: []string{"f"},
				Usage:   "是否过滤掉 github.com/xiaoxuan6/go-package-example 中已存在的库",
			},
		},
		Action: func(context *cli.Context) error {
			account := context.String("account")
			accounts := strings.Split(account, "|")

			var (
				wg       sync.WaitGroup
				userChan = make(chan map[string][]*github2.Repository)
			)
			for _, username := range accounts {
				wg.Add(1)
				go func(username string) {
					defer wg.Done()

					if err := verifyUser(username); err != nil {
						userChan <- map[string][]*github2.Repository{username: []*github2.Repository{}}
						return
					}

					starAllRepos := downloadRepos(username)
					userChan <- map[string][]*github2.Repository{username: starAllRepos}
				}(username)
			}

			wg.Add(1)
			go func(filter bool) {
				defer wg.Done()
				if filter == false {
					return
				}

				fetchHistory()
			}(context.Bool("filter"))

			go func() {
				wg.Wait()
				close(userChan)
			}()

			for user := range userChan {
				for username, starRepos := range user {
					if len(starRepos) > 0 {
						fmt.Printf("fetch %s star repository success by start translate description.\n", username)
						wg.Add(1)
						go translateDescription(&wg, username, starRepos)
					} else {
						userStarRepos[username] = []*global.GRepository{
							&global.GRepository{
								Repository: &github2.Repository{},
							},
						}
						fmt.Println(fmt.Sprintf("username %s download done.", username))
					}
				}
			}

			wg.Wait()
			saveFile()
			fmt.Println("all download done.")
			return nil
		},
	}
}

func verifyUser(account string) error {
	users := github.SearchOwner(account)
	if len(users) == 0 {
		return errors.New(fmt.Sprintf("未找到 [%s] 用户！", account))
	}

	target := false
	for _, user := range users {
		if strings.Compare(account, user.GetLogin()) == 0 {
			target = true
		}
	}

	if target == false {
		return errors.New(fmt.Sprintf("未找到 [%s] 用户！", account))
	}

	return nil
}

func downloadRepos(username string) []*github2.Repository {
	page, starAllRepos := 1, make([]*github2.Repository, 0)
	for page != 0 {
		starRepos, nextPage := services.FetchStarRepos(username, page, 100)
		if len(starRepos) > 0 {
			starAllRepos = append(starAllRepos, starRepos...)
		}
		page = nextPage
	}

	return starAllRepos
}

var (
	userStarRepos = make(map[string][]*global.GRepository)
)

func translateDescription(wg *sync.WaitGroup, username string, starRepos []*github2.Repository) {
	defer wg.Done()

	var swg sync.WaitGroup
	gStarRepos := make([]*global.GRepository, 0)
	for _, repos := range starRepos {
		repos := repos
		if slices.Contains(history, repos.GetFullName()) {
			fmt.Println(fmt.Sprintf("repo %s exists", repos.GetFullName()))
			continue
		} else {
			swg.Add(1)
			go func() {
				defer swg.Done()

				desc, stat := translate.Translation(repos.GetDescription())
				gStarRepos = append(gStarRepos, &global.GRepository{
					Repository:           repos,
					DescriptionTranslate: desc,
					TranslateStat:        stat,
				})
			}()
		}
	}

	swg.Wait()
	userStarRepos[username] = gStarRepos
}

func saveFile() {
	f := excelize.NewFile()

	for username, repos := range userStarRepos {
		index, err := f.NewSheet(username)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		header := []interface{}{
			"Repository",
			"Description",
			"Chinese Description",
			"Stars",
			"Forks",
			"Language",
		}
		_ = f.SetSheetRow(username, fmt.Sprintf("A1"), &header)

		for i, repo := range repos {
			r := repo.Repository

			row := []interface{}{
				fmt.Sprintf("https://github.com/%s", r.GetFullName()),
				r.GetDescription(),
				repo.DescriptionTranslate,
				r.GetStargazersCount(),
				r.GetForksCount(),
				r.GetLanguage(),
			}
			if err = f.SetSheetRow(username, fmt.Sprintf("A%d", i+2), &row); err != nil {
				fmt.Println(err)
			}
		}

		f.SetActiveSheet(index)
	}

	_ = f.DeleteSheet("Sheet1")
	_ = f.SaveAs("star_repos.xlsx")
}

var history = make([]string, 100)

func fetchHistory() {
	urls := []string{
		"https://github-mirror.us.kg/https:/github.com/xiaoxuan6/go-package-example/blob/main/README_PHP.md",
		"https://github-mirror.us.kg/https:/github.com/xiaoxuan6/go-package-example/blob/main/README_OTHER.md",
		"https://github-mirror.us.kg/https:/github.com/xiaoxuan6/go-package-example/blob/main/README.md",
	}

	var (
		wg     sync.WaitGroup
		client = http.Client{
			Timeout: 3 * time.Second,
		}
	)

	for _, url := range urls {
		wg.Add(1)

		url := url
		go func() {
			defer wg.Done()
			response, err := client.Get(url)
			if err != nil {
				return
			}

			defer response.Body.Close()
			f := bufio.NewReader(response.Body)
			for {
				line, _, err := f.ReadLine()
				if err == io.EOF {
					break
				}

				x := xurls.Relaxed()
				domain := x.FindString(string(line))
				domain = strings.ReplaceAll(domain, "github.com/", "")
				if len(domain) > 1 {
					history = append(history, domain)
				}
			}
		}()
	}

	wg.Wait()
}
