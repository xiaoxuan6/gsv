package tui

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	github2 "github.com/google/go-github/v64/github"
	"github.com/samber/lo"
	"github.com/skratchdot/open-golang/open"
	"github.com/xiaoxuan6/gsv/pkg/github"
	"github.com/xiaoxuan6/gsv/pkg/global"
	"github.com/xiaoxuan6/gsv/pkg/spinner"
	"github.com/xiaoxuan6/gsv/pkg/translate"
	"github.com/xiaoxuan6/gsv/services"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var regex = `【(.*?)】（(.*?)）`

func RenderCurrentList() {
	currentStarRepos := global.AccountsAllStarRepos[global.CurrentAccount]
	items := services.CheckItem(currentStarRepos)
	nextPage := global.AccountsStarReposNextPage[global.CurrentAccount]
	RenderList(items, nextPage, len(items))
}

func RenderList(items []string, page, total int) {
	items = append(items, fmt.Sprintf("【footer】（%s） - %s", strconv.Itoa(page), strconv.Itoa(total)))
	items = lo.Map(items, func(item string, index int) string {
		return fmt.Sprintf("%d、%s", index+1, item)
	})

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "All stars repos"
	l.Rows = items
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	l.TextStyle = ui.NewStyle(ui.ColorWhite)
	l.SetRect(0, 5, 200, 30)
	ui.Render(ListHelp(), l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "a":
			ui.Clear()
			ui.Close()
			RenderAccounts()
			os.Exit(0)
		case "<C-c>":
			ui.Clear()
			ui.Close()
			os.Exit(0)
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "b":
			l.ScrollBottom()
		case "t":
			l.ScrollTop()
		case "s":
			ui.Clear()
			ui.Close()
			RenderSearch()
			os.Exit(0)
		case "<Enter>":
			ui.Clear()
			ui.Close()

			str := l.Rows[l.SelectedRow]
			re1 := regexp.MustCompile(regex).FindAllStringSubmatch(str, -1)
			if len(re1) != 1 {
				fmt.Println("选中数据无效")
				os.Exit(0)
			}

			item := re1[0][2]
			if len(item) < 1 && len(re1[0][1]) < 1 {
				fmt.Println("github repos empty")
				os.Exit(0)
			}

			i, _ := strconv.Atoi(item)
			if i != 0 {
				ReloadRenderList(i)
			} else {
				RenderRepos(re1[0][1])
			}
			os.Exit(0)
		}
		ui.Render(l)
	}
}

func RenderSearch() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = "请输入 GITHUB 账号"
	p.Text = "直接输入然后回车......"
	p.PaddingTop = 1
	p.SetRect(0, 0, 50, 5)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	ui.Render(p, SearchHelp())

	input := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			ui.Clear()
			ui.Close()
			os.Exit(0)
			return
		case "c":
			ui.Clear()
			ui.Close()
			RenderCurrentList()
			os.Exit(0)
		case "<Enter>":
			ui.Clear()
			ui.Close()

			if len(p.Text) < 1 {
				RenderSearch()
			} else {
				fetchRepos(p.Text)
			}
			os.Exit(0)
		case "<C-r>":
			ui.Clear()
			ui.Close()
			RenderSearch()
			os.Exit(0)
		default:
			if e.Type == ui.KeyboardEvent {
				input += e.ID
				p.Text = input
				ui.Render(p)
			}
		}
	}
}

func fetchRepos(username string) {
	if _, ok := global.AccountsAllStarRepos[username]; ok {
		global.CurrentAccount = username
		RenderCurrentList()
	} else {
		owners := spinner.RunF[[]string]("fetching github owners ", func() []string {
			users := github.SearchOwner(username)
			owners := make([]string, 0)
			for _, value := range users {
				owners = append(owners, value.GetLogin())
			}

			return owners
		})

		if len(owners) == 0 {
			owners = append(owners, "暂无数据")
		}

		if err := ui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		defer ui.Close()

		l := widgets.NewList()
		l.Title = "All owners"
		l.Rows = owners
		l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
		l.TextStyle = ui.NewStyle(ui.ColorWhite)
		l.SetRect(0, 0, 50, 10)
		ui.Render(l, FetchReposHelp())

		uiEvents := ui.PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "<C-c>":
				ui.Clear()
				ui.Close()
				os.Exit(0)
				return
			case "j", "<Down>":
				l.ScrollDown()
			case "k", "<Up>":
				l.ScrollUp()
			case "s":
				ui.Clear()
				ui.Close()
				RenderSearch()
				os.Exit(0)
			case "e", "<Enter>":
				ui.Clear()
				ui.Close()

				owner := l.Rows[l.SelectedRow]
				if strings.Compare(owner, "暂无数据") == 0 {
					RenderSearch()
				} else {
					global.CurrentAccount = owner
					ReloadRenderList(1)
				}
				os.Exit(0)
			}
			ui.Render(l)
		}
	}
}

func ReloadRenderList(page int) {
	allRepos, nextPage := spinner.RunF2[[]string, int]("fetching github stars repos data ", func() ([]string, int) {
		return services.FetchDataWithPage(global.CurrentAccount, page)
	})
	RenderList(allRepos, nextPage, len(allRepos))
}

func RenderRepos(repos string) {
	accountStarRepos := global.AccountsAllStarRepos[global.CurrentAccount]

	target := false
	for _, val := range accountStarRepos {
		if strings.Compare(*val.Repository.FullName, repos) == 0 {
			target = true
			RenderTable(val.Repository, val.DescriptionZh)
		}
	}

	if target == false {
		fmt.Printf("暂无 [%s] 数据！", repos)
		return
	}
}

func RenderTable(repos *github2.Repository, description string) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	table.Rows = [][]string{
		{"repos", "desc", "language", "stars", "forks"},
		{repos.GetFullName(), description, repos.GetLanguage(), strconv.Itoa(repos.GetStargazersCount()), strconv.Itoa(repos.GetForksCount())},
	}
	table.ColumnWidths = []int{30, 80, 10, 10, 15}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.TextAlignment = ui.AlignCenter
	table.SetRect(0, 0, 145, 10)
	ui.Render(table, TableHelp())

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "a":
			ui.Clear()
			ui.Close()
			RenderAccounts()
			os.Exit(0)
		case "<C-c>":
			ui.Clear()
			ui.Close()
			os.Exit(0)
			return
		case "s":
			ui.Clear()
			ui.Close()
			RenderSearch()
			os.Exit(0)
		case "t":
			ui.Clear()
			ui.Close()

			desc := spinner.RunF[string]("translate doing", func() string {
				return translate.Translation(description)
			})

			RenderTable(repos, desc)
			os.Exit(0)
		case "d":
			ui.Clear()
			ui.Close()

			allStarsRepos := global.AccountsAllStarRepos[global.CurrentAccount]
			allStarsRepos = lo.FilterMap(allStarsRepos, func(item *global.GRepository, _ int) (*global.GRepository, bool) {
				if strings.Compare(*item.Repository.FullName, repos.GetFullName()) == 0 {
					return item, false
				}
				return item, true
			})
			global.AccountsAllStarRepos[global.CurrentAccount] = allStarsRepos

			items := services.CheckItem(allStarsRepos)
			RenderList(items, global.AccountsStarReposNextPage[global.CurrentAccount], len(items))
			os.Exit(0)
		case "o":
			_ = open.Run(fmt.Sprintf("https://github.com/%s", repos.GetFullName()))
		case "r":
			ui.Clear()
			ui.Close()

			RenderCurrentList()
			os.Exit(0)
		}
		ui.Render(table, TableHelp())
	}
}

func RenderAccounts() {
	accounts := lo.Keys(global.AccountsStarReposNextPage)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "All account"
	l.Rows = accounts
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	l.TextStyle = ui.NewStyle(ui.ColorWhite)
	l.SetRect(0, 5, 30, 30)
	ui.Render(CurrentAccount(), l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			ui.Clear()
			ui.Close()
			os.Exit(0)
		case "c":
			ui.Clear()
			ui.Close()
			RenderCurrentList()
			os.Exit(0)
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "s":
			ui.Clear()
			ui.Close()
			RenderSearch()
			os.Exit(0)
		case "<Enter>":
			ui.Clear()
			ui.Close()

			username := l.Rows[l.SelectedRow]
			fetchRepos(strings.TrimSpace(username))
			os.Exit(0)
		}
		ui.Render(l)
	}
}
