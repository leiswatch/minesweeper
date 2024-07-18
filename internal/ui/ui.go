package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	greenColor = lipgloss.Color("#a6e3a1")
	blueColor  = lipgloss.Color("#89b4fa")
	redColor   = lipgloss.Color("#f38ba8")
	whiteColor = lipgloss.Color("#cdd6f4")
	grayColor  = lipgloss.Color("#585b70")
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
				Foreground(grayColor),
		).Rows(rows...)
}

func newTextStyle(bold bool) lipgloss.Style {
	return lipgloss.NewStyle().Height(1).Bold(bold)
}

func RenderBoldText(s string) string {
	return newTextStyle(true).Render(s)
}

func RenderGreenText(s string) string {
	return newTextStyle(true).Foreground(greenColor).Render(s)
}

func RenderRedText(s string) string {
	return newTextStyle(true).Foreground(redColor).Render(s)
}

func RenderBlueText(s string) string {
	return newTextStyle(true).Foreground(blueColor).Render(s)
}

func RenderWhiteText(s string) string {
	return newTextStyle(false).Foreground(whiteColor).Render(s)
}

func RenderGrayText(s string) string {
	return newTextStyle(false).Foreground(grayColor).Render(s)
}

