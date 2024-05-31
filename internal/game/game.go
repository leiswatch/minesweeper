package game

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type position struct {
	x int
	y int
}

type gameModel struct {
	board    [][]byte
	position position
}

func newBoard(width, height int) [][]byte {
	board := make([][]byte, width)

	for i := 0; i < width; i++ {
		board[i] = make([]byte, height)
	}

	for i := range board {
		for j := range board[i] {
			board[i][j] = '.'
		}
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
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m gameModel) View() string {
	// The header
	s := "Minesweeper\n\n"

	var cursor byte = 'x'

	// Iterate over our choices
	for i := range len(m.board) {
		for j := range len(m.board[i]) {
			if m.position.x == i && m.position.y == j {
				m.board[i][j] = cursor
			} else {
				m.board[i][j] = '.'
			}
		}

		// Render the row
		s += fmt.Sprintf("%s\n", m.board[i])
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func StartGame() {
	p := tea.NewProgram(newGameModel(10, 15))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
