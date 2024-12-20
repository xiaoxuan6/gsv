package commands

import (
	"github.com/briandowns/spinner"
	"github.com/urfave/cli/v2"
	"github.com/xiaoxuan6/gsv/pkg/global"
	"github.com/xiaoxuan6/gsv/pkg/tui"
	"github.com/xiaoxuan6/gsv/services"
	"time"
)

func Search() *cli.Command {
	return &cli.Command{
		Name:        "search",
		Usage:       "查找指定用户所有的 star repos",
		Description: "查找指定用户所有的 star repos（默认 100 条）",
		Aliases:     []string{"s"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "account",
				Aliases:  []string{"a"},
				Required: true,
				Usage:    "需要查找的 github 账号",
			},
			&cli.IntFlag{
				Name:    "page",
				Aliases: []string{"p"},
				Value:   1,
				Usage:   "需要查询的页数",
			},
			&cli.IntFlag{
				Name:    "pageCount",
				Aliases: []string{"pc"},
				Value:   100,
				Usage:   "每页查询的个数",
			},
			&cli.BoolFlag{
				Name:    "filter",
				Aliases: []string{"f"},
				Value:   false,
				Usage:   "是否过滤掉 github.com/xiaoxuan6/go-package-example 中已存在的库",
			},
		},
		Action: func(c *cli.Context) error {
			global.FilterAction = c.Bool("filter")
			if global.FilterAction {
				go services.FetchHistory()
			}

			s := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
			s.Prefix = "fetching github stars repos data "
			s.FinalMSG = "done"
			s.Start()

			global.CurrentAccount = c.String("account")
			global.PageCount = c.Int("pageCount")
			allRepos, nextPage := services.FetchData(c.String("account"), c.Int("page"), global.PageCount)
			global.AccountsStarReposNextPage[global.CurrentAccount] = nextPage
			s.Stop()

			tui.RenderList(allRepos, nextPage, len(allRepos))
			return nil
		},
	}
}
