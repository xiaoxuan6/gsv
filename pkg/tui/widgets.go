package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func ListHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、j：上一条、K：下一条、b：置低、t：置顶、s：搜索 github owner、e|Enter：进入详情页"
	p.Border = false
	p.SetRect(0, 5, 110, 1)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}

func SearchHelp() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "Ctrl+c：终止程序、Ctrl+r：重新搜索、Ctrl+e|Enter：进入下一步"
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
	p.Text = "Ctrl+c：终止程序 、s：重新搜索、r：github owner stars list、o：使用浏览器打开链接、t：翻译 desc、d:删除当前数据"
	p.Border = false
	p.SetRect(0, 20, 120, 10)
	p.TextStyle.Fg = ui.ColorMagenta

	return p
}
