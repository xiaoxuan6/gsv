package commands

import (
	"github.com/briandowns/spinner"
	"github.com/urfave/cli/v2"
	"github.com/xiaoxuan6/gsv/pkg/global"
	"github.com/xiaoxuan6/gsv/pkg/tui"
	"github.com/xiaoxuan6/gsv/services"
	"time"
)

func AllRepos() *cli.Command {
	return &cli.Command{
		Name:        "all-repos",
		Usage:       "查找指定用户所有的 star repos",
		Description: "查找指定用户所有的 star repos（默认 100 条）",
		Aliases:     []string{"all"},
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
		},
		Action: func(c *cli.Context) error {
			s := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
			s.Prefix = "fetching github stars repos data "
			s.FinalMSG = "done"
			s.Start()

			global.CurrentAccount = c.String("account")
			allRepos, nextPage := services.FetchData(c.String("account"), c.Int("page"), c.Int("pageCount"))
			global.AccountsStarReposNextPage[global.CurrentAccount] = nextPage
			s.Stop()

			tui.RenderList(allRepos, nextPage, len(allRepos))
			return nil
		},
	}
}
