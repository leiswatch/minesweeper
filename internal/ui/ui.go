package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func NewBoardTable(rows [][]string) *table.Table {
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
				Foreground(lipgloss.Color("#585b70")),
		).Rows(rows...)
}
