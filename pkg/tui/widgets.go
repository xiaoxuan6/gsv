package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/xiaoxuan6/gsv/pkg/global"
)

func ListHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、Ctrl+t: 翻译当前页所有数据、Enter：进入详情页\na：显示所有用户、b：置低、c: 返回列表、j：上一条、K：下一条、l: 显示所有语言、o: 打开链接、s：搜索 github 用户、t：置顶"
	p.Border = false
	p.SetRect(0, 5, 160, 1)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func LanguageListHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、Enter：进入详情页、a：显示所有用户、b：置低、c: 返回列表、j：上一条、K：下一条、t：置顶"
	p.Border = false
	p.SetRect(0, 5, 140, 1)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func SearchHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、Ctrl+r：重新搜索、Enter：进入下一步、tab: 取消搜索"
	p.Border = false
	p.SetRect(0, 5, 75, 10)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func FetchReposHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、Enter：展示 github stars、j：上一条、K：下一条、s：重新搜索"
	p.Border = false
	p.SetRect(0, 20, 95, 10)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func TableHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、a：显示所有的用户、o：使用浏览器打开链接、r：github stars 列表、s：重新搜索、t：翻译 desc"
	p.Border = false
	p.SetRect(0, 30, 140, 15)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func TableDesc(desc string) *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Title = " description "
	p.Border = true
	p.Text = desc
	p.PaddingTop = 1
	p.SetRect(0, 10, 165, 15)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorWhite

	return p
}

func CurrentAccount() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "当前用户：" + global.CurrentAccount + "\n\nCtrl+c：终止程序、Enter：进入详情页、c: 取消、j：上一条、K：下一条、s：重新搜索"
	p.Border = false
	p.SetRect(0, 0, 100, 5)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}
