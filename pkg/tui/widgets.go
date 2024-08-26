package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/xiaoxuan6/gsv/pkg/global"
)

func ListHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、a：显示所有的用户、j：上一条、K：下一条、b：置低、t：置顶、s：搜索 github 用户、Enter：进入详情页"
	p.Border = false
	p.SetRect(0, 5, 120, 1)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func SearchHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、Ctrl+r：重新搜索、c: 取消搜索、Enter：进入下一步"
	p.Border = false
	p.SetRect(0, 5, 75, 10)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func FetchReposHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、j：上一条、K：下一条、s：重新搜索、e|Enter：展示 github owner stars"
	p.Border = false
	p.SetRect(0, 20, 95, 10)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func TableHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序 、a：显示所有的用户、s：重新搜索、r：github stars 列表、o：使用浏览器打开链接、t：翻译 desc、d:删除当前数据"
	p.Border = false
	p.SetRect(0, 20, 140, 10)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func CurrentAccount() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "当前用户：" + global.CurrentAccount + "\n\nCtrl+c：终止程序 、c: 取消、s：重新搜索、j：上一条、K：下一条、Enter：进入详情页"
	p.Border = false
	p.SetRect(0, 0, 100, 5)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}
