package game

import (
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type state struct {
	win  bool
	lose bool
}

type square struct {
	isFlag bool
	isMine bool
	el     string
}

type board = [][]square

type cursorPos struct {
	x int
	y int
}

type model struct {
	board     board
	cursorPos cursorPos
	state     state
}

const (
	cursorEl string = "X"
	flagEl   string = "F"
	hiddenEl string = " "
	mineEl   string = "B"
)

type bombPosition struct {
	x int
	y int
}

func generateBombPositions(n int) []bombPosition {
	seed := uint64(time.Now().UnixNano())
	r := rand.New(rand.NewPCG(seed, seed))

	m := make(map[bombPosition]int)

	for len(m) < n {
		x := r.IntN(9)
		y := r.IntN(9)

		bp := bombPosition{
			x,
			y,
		}

		_, ok := m[bp]

		if ok {
			m[bp] = m[bp] + 1
		} else {
			m[bp] = 1
		}
	}

	bombPositions := make([]bombPosition, 0, n)

	for k := range m {
		bombPositions = append(bombPositions, k)
	}

	return bombPositions
}

func newBoard(width, height int) board {
	board := make(board, width)
	bombPositions := generateBombPositions(10)

	for i := range width {
		board[i] = make([]square, height)

		for j := range height {
			board[i][j] = square{
				isFlag: false,
				isMine: slices.Contains(bombPositions, bombPosition{
					x: i,
					y: j,
				}),
				el: hiddenEl,
			}
		}
	}

	return board
}

func newModel(width, height int) model {
	return model{
		board: newBoard(width, height),
		cursorPos: cursorPos{
			x: width - 1,
			y: height - 1,
		},
	}
}

func (m model) Init() tea.Cmd {
	// Just returno `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if !(m.state.win || m.state.lose) && m.cursorPos.y > 0 {
				m.cursorPos.y--
			}

		// The "down" and "j" keys move the cursor down
		case "left", "l":
			if !(m.state.win || m.state.lose) && m.cursorPos.y < len(m.board[0])-1 {
				m.cursorPos.y++
			}

		case "down", "j":
			if !(m.state.win || m.state.lose) && m.cursorPos.x < len(m.board)-1 {
				m.cursorPos.x++
			}

		case "up", "k":
			if !(m.state.win || m.state.lose) && m.cursorPos.x > 0 {
				m.cursorPos.x--
			}

		case "f":
			m.board[m.cursorPos.x][m.cursorPos.y].isFlag = !m.board[m.cursorPos.x][m.cursorPos.y].isFlag

		case "enter":
			if m.board[m.cursorPos.x][m.cursorPos.y].isMine && !m.board[m.cursorPos.x][m.cursorPos.y].isFlag {
				m.state.lose = !m.state.lose
			}

		case "r":
			if m.state.lose || m.state.win {
				m.board = newBoard(9, 9)
				m.state.lose = false
				m.state.win = false
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Minesweeper\n\n"

	for i := range len(m.board) {
		for j := range len(m.board[i]) {
			if m.board[i][j].isMine && m.state.lose {
				m.board[i][j].el = mineEl
			} else if m.cursorPos.x == i && m.cursorPos.y == j {
				m.board[i][j].el = cursorEl
			} else if m.board[i][j].isFlag {
				m.board[i][j].el = flagEl
			} else {
				m.board[i][j].el = hiddenEl
			}
		}
	}

	rows := make([][]string, len(m.board))

	for i := range len(m.board) {
		rows[i] = make([]string, len(m.board[i]))

		for j := range len(m.board[i]) {
			rows[i][j] = m.board[i][j].el
		}
	}

	boardTable := createBoardTable(rows)

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
	p := tea.NewProgram(newModel(9, 9))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
