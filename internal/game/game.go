package game

import (
	"fmt"
	"os"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type position struct {
	x int
	y int
}

type gameModel struct {
	board    [][]string
	position position
	flags    []position
}

const (
	cursor string = "x"
	flag   string = "F"
	hidden string = "."
)

func newBoard(width, height int) [][]string {
	board := make([][]string, width)

	for i := range width {
		board[i] = make([]string, height)
	}

	return board
}

func newGameModel(width, height int) gameModel {
	return gameModel{
		board: newBoard(width, height),
		position: position{
			x: width - 1,
			y: height - 1,
		},
	}
}

func (m gameModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "right", "h":
			if m.position.y > 0 {
				m.position.y--
			}

		// The "down" and "j" keys move the cursor down
		case "left", "l":
			if m.position.y < len(m.board[0])-1 {
				m.position.y++
			}

		case "down", "j":
			if m.position.x < len(m.board)-1 {
				m.position.x++
			}

		case "up", "k":
			if m.position.x > 0 {
				m.position.x--
			}

		case "f":
			currPosition := position{
				x: m.position.x,
				y: m.position.y,
			}

			if !slices.Contains(m.flags, currPosition) {
				m.flags = append(m.flags, currPosition)
			} else {
				idx := slices.Index(m.flags, currPosition)
				m.flags = slices.Delete(m.flags, idx, idx+1)
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m gameModel) View() string {
	// The header
	s := "Minesweeper\n\n"

	for i := range len(m.board) {
		for j := range len(m.board[i]) {
			currCell := position{
				x: i,
				y: j,
			}

			if m.position.x == currCell.x && m.position.y == currCell.y {
				m.board[currCell.x][currCell.y] = cursor
			} else if slices.Contains(m.flags, currCell) {
				m.board[currCell.x][currCell.y] = flag
			} else {
				m.board[currCell.x][currCell.y] = hidden
			}
		}
	}

	boardTable := createBoardTable(m.board)

	s += fmt.Sprintf("%s\n", boardTable)
	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func createBoardTable(rows [][]string) *table.Table {
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

func StartGame() {
	p := tea.NewProgram(newGameModel(8, 8))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
