package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	BlueColor  = lipgloss.Color("#89b4fa")
	GrayColor  = lipgloss.Color("#585b70")
	GreenColor = lipgloss.Color("#a6e3a1")
	PinkColor  = lipgloss.Color("#f5c2e7")
	RedColor   = lipgloss.Color("#f38ba8")
	WhiteColor = lipgloss.Color("#cdd6f4")
)

func NewBoard(rows [][]string) *table.Table {
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		StyleFunc(
			func(row, col int) lipgloss.Style {
				return lipgloss.NewStyle().
					Align(
						lipgloss.Position(lipgloss.Center),
						lipgloss.Position(lipgloss.Center),
					).
					PaddingLeft(1).
					PaddingRight(1)
			}).
		BorderStyle(
			lipgloss.NewStyle().
				Foreground(GrayColor),
		).Rows(rows...)
}

func newTextStyle(bold bool) lipgloss.Style {
	return lipgloss.NewStyle().Height(1).Bold(bold)
}

func RenderText(s string, color lipgloss.Color, bold bool) string {
	return newTextStyle(bold).Foreground(color).Render(s)
}
