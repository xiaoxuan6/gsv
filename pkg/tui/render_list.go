package tui

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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
	"sync"
	"time"
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

	title := "All stars repos"
	if len(currentLanguage) > 1 {
		title = fmt.Sprintf("Language【%s】stars repos", currentLanguage)
		currentLanguage = ""
		global.SelectedRow = 0
	}

	l := widgets.NewList()
	l.Title = title
	l.Rows = items
	l.SelectedRow = global.SelectedRow
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
			global.SelectedRow = l.SelectedRow
			RenderAccounts()
			os.Exit(0)
		case "<C-c>":
			ui.Clear()
			ui.Close()
			os.Exit(0)
			return
		case "<C-t>":
			ui.Clear()
			ui.Close()

			var wg sync.WaitGroup
			newStarRepos := make([]*global.GRepository, 0)
			allStarRepos := global.AccountsAllStarRepos[global.CurrentAccount]
			for _, repo := range allStarRepos {
				if repo.TranslateStat == false {
					wg.Add(1)
					repo := repo
					go func() {
						defer wg.Done()
						desc, stat := translate.Translation(repo.Repository.GetDescription())
						repo.TranslateStat = stat
						repo.DescriptionTranslate = desc
						repo.DescriptionZh = fmt.Sprintf("【%s】（%s） - %s", repo.Repository.GetFullName(), desc, repo.Repository.GetLanguage())
						newStarRepos = append(newStarRepos, repo)
					}()
				} else {
					newStarRepos = append(newStarRepos, repo)
				}
			}
			wg.Wait()

			global.AccountsAllStarRepos[global.CurrentAccount] = newStarRepos
			go services.LanguageCategory(newStarRepos, global.AccountsAllLanguages[global.CurrentAccount])
			RenderCurrentList()
			os.Exit(0)
			return
		case "c":
			ui.Clear()
			ui.Close()
			RenderCurrentList()
			os.Exit(0)
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
			global.SelectedRow = l.SelectedRow
			RenderSearch()
			os.Exit(0)
		case "l":
			ui.Clear()
			ui.Close()
			global.SelectedRow = l.SelectedRow
			RenderLanguagesList()
			os.Exit(0)
		case "o":
			str := l.Rows[l.SelectedRow]
			reg := regexp.MustCompile(`【(.*?)】（`).FindStringSubmatch(str)
			if len(reg) < 2 {
				fmt.Printf("系统错误！")
				os.Exit(0)
			}

			_ = open.Run(fmt.Sprintf("https://github.com/%s", reg[1]))
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
				global.SelectedRow = 0
				ReloadRenderList(i)
			} else {
				global.SelectedRow = l.SelectedRow
				RenderRepos(re1[0][1])
			}
			os.Exit(0)
		}
		ui.Render(ListHelp(), l)
	}
}

var currentLanguage string

func RenderLanguagesList() {
ITEM:
	items := global.AccountsAllLanguages[global.CurrentAccount]
	if len(items) < 1 {
		spinner.RunF0("sleep", func() {
			time.Sleep(3)
		})

		goto ITEM
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "All language list"
	l.Rows = items
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	l.TextStyle = ui.NewStyle(ui.ColorWhite)
	l.SetRect(0, 5, 30, 20)
	ui.Render(LanguageListHelp(), l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			ui.Clear()
			ui.Close()
			os.Exit(0)
		case "a":
			ui.Clear()
			ui.Close()
			RenderAccounts()
			os.Exit(0)
		case "c":
			ui.Clear()
			ui.Close()
			RenderCurrentList()
			os.Exit(0)
		case "j":
			l.ScrollDown()
		case "k":
			l.ScrollUp()
		case "b":
			l.ScrollBottom()
		case "t":
			l.ScrollTop()
		case "<Enter>":
			ui.Clear()
			ui.Close()

			language := l.Rows[l.SelectedRow]
			if languageStarRepos, ok := global.AccountsLanguageStarRepose[global.CurrentAccount][language]; ok {
				currentLanguage = language

				languageItems := services.CheckItem(languageStarRepos)
				nextPage := global.AccountsStarReposNextPage[global.CurrentAccount]
				RenderList(languageItems, nextPage, len(languageItems))
			}
			os.Exit(0)
		}
		ui.Render(LanguageListHelp(), l)
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
			case "<Enter>":
				ui.Clear()
				ui.Close()

				owner := l.Rows[l.SelectedRow]
				if strings.Compare(owner, "暂无数据") == 0 {
					RenderSearch()
				} else {
					global.CurrentAccount = owner
					global.SelectedRow = 0
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
			RenderTable(val, val.DescriptionTranslate)
		}
	}

	if target == false {
		fmt.Printf("暂无 [%s] 数据！\n当前用户：%s, 总数据：%d", repos, global.CurrentAccount, len(accountStarRepos))
		return
	}
}

func RenderTable(gRepos *global.GRepository, description string) {
	repos := gRepos.Repository
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	table.Rows = [][]string{
		{"id", "repos", "language", "stars", "forks"},
		{strconv.Itoa(global.SelectedRow + 1), repos.GetFullName(), repos.GetLanguage(), strconv.Itoa(repos.GetStargazersCount()), strconv.Itoa(repos.GetForksCount())},
	}
	table.ColumnWidths = []int{10, 30, 10, 10, 20}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.TextAlignment = ui.AlignCenter
	table.SetRect(0, 0, 80, 10)
	ui.Render(table, TableDesc(description), TableHelp())

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

			var (
				desc string
				stat bool
			)
			if gRepos.TranslateStat != true {
				desc, stat = spinner.RunF2[string, bool]("translate doing", func() (string, bool) {
					result, ok := translate.Translation(description)
					result = strings.ReplaceAll(strings.ReplaceAll(result, " | ", ""), "|", "")
					return result, ok
				})

				go func() {
					allStarRepos := global.AccountsAllStarRepos[global.CurrentAccount]
					allStarRepos = lo.FilterMap(allStarRepos, func(item *global.GRepository, _ int) (*global.GRepository, bool) {
						if strings.Compare(item.Repository.GetFullName(), repos.GetFullName()) == 0 {
							return &global.GRepository{
								Repository:           item.Repository,
								DescriptionZh:        fmt.Sprintf("【%s】（%s） - %s", repos.GetFullName(), desc, repos.GetLanguage()),
								DescriptionTranslate: desc,
								TranslateStat:        stat,
							}, true
						}

						return item, true
					})

					global.AccountsAllStarRepos[global.CurrentAccount] = allStarRepos
				}()
			} else {
				desc = gRepos.DescriptionTranslate
			}

			RenderTable(gRepos, desc)
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
		ui.Render(table, TableDesc(description), TableHelp())
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
			if strings.Compare(username, global.CurrentAccount) != 0 {
				global.SelectedRow = 0
			}
			fetchRepos(strings.TrimSpace(username))
			os.Exit(0)
		}
		ui.Render(l)
	}
}
