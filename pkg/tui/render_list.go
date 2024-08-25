package tui

import (
	"fmt"
	"github.com/briandowns/spinner"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	github2 "github.com/google/go-github/v64/github"
	"github.com/samber/lo"
	"github.com/xiaoxuan6/gsv/pkg/github"
	"github.com/xiaoxuan6/gsv/pkg/global"
	"github.com/xiaoxuan6/gsv/services"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var regex = `【(.*?)】（(.*?)）`

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
		case "e", "<Enter>":
			ui.Clear()
			ui.Close()

			str := l.Rows[l.SelectedRow]
			re1 := regexp.MustCompile(regex).FindAllStringSubmatch(str, -1)
			if len(re1) != 1 {
				fmt.Println("选中数据无效")
				os.Exit(0)
			}

			item := re1[0][2]
			if len(item) < 1 {
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
		case "<C-e>", "<Enter>":
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
	if allStarRepos, ok := global.AccountsAllStarRepos[username]; ok {
		items := services.CheckItem(allStarRepos)
		RenderList(items, 0, len(items))
	} else {
		s := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
		s.Prefix = "fetching github owners "
		s.FinalMSG = "done"
		s.Start()

		users := github.SearchOwner(username)
		owners := make([]string, 0)
		for _, value := range users {
			owners = append(owners, value.GetLogin())
		}

		s.Stop()
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
	s := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
	s.Prefix = "fetching github stars repos data "
	s.FinalMSG = "done"
	s.Start()

	allRepos, nextPage := services.FetchDataWithPage(global.CurrentAccount, page)

	s.Stop()
	RenderList(allRepos, nextPage, len(allRepos))
}

func RenderRepos(repos string) {
	accountStarRepos := global.AccountsAllStarRepos[global.CurrentAccount]

	target := false
	for _, val := range accountStarRepos {
		if strings.Compare(*val.Repository.FullName, repos) == 0 {
			target = true
			RenderTable(val.Repository)
		}
	}

	if target == false {
		fmt.Printf("暂无 [%s] 数据！", repos)
		return
	}
}

func RenderTable(repos *github2.Repository) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	table.Title = repos.GetFullName()
	table.Rows = [][]string{
		{"repos", "desc", "language", "stars", "forks"},
		{repos.GetFullName(), repos.GetDescription(), repos.GetLanguage(), strconv.Itoa(repos.GetStargazersCount()), strconv.Itoa(repos.GetForksCount())},
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
		case "r":
			ui.Clear()
			ui.Close()

			s := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
			s.Prefix = "fetching github stars repos data "
			s.FinalMSG = "done"
			s.Start()

			currentRepos := global.AccountsAllStarRepos[global.CurrentAccount]
			items := services.CheckItem(currentRepos)

			s.Stop()

			RenderList(items, 0, len(items))
			os.Exit(0)
		}
	}
}
