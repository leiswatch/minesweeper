package game

import (
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/leiswatch/minesweeper/internal/ui"
)

type state struct {
	win  bool
	lose bool
}

type square struct {
	el     string
	isFlag bool
	isMine bool
	count  int
}

type board = [][]square

type coordinates struct {
	x int
	y int
}

type model struct {
	board     board
	cursorPos coordinates
	state     state
}

const (
	cursorEl string = "X"
	flagEl   string = "F"
	hiddenEl string = " "
	mineEl   string = "B"
)

var directions = [8]coordinates{
	{x: 1, y: 0},
	{x: -1, y: 0},
	{x: 0, y: 1},
	{x: 0, y: -1},
	{x: 1, y: -1},
	{x: 1, y: 1},
	{x: -1, y: 1},
	{x: -1, y: 1},
}

func generateBombPositions(n int) []coordinates {
	seed := uint64(time.Now().UnixNano())
	r := rand.New(rand.NewPCG(seed, seed))

	m := make(map[coordinates]int)

	for len(m) < n {
		x := r.IntN(9)
		y := r.IntN(9)

		bp := coordinates{
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

	bombPositions := make([]coordinates, 0, n)

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
				isMine: slices.Contains(bombPositions, coordinates{
					x: i,
					y: j,
				}),
				el:    hiddenEl,
				count: -1,
			}
		}
	}

	return board
}

func newModel(width, height int) model {
	return model{
		board: newBoard(width, height),
		cursorPos: coordinates{
			x: width - 1,
			y: height - 1,
		},
	}
}

func revealCells(b board, row, col int) {
	if row < 0 || col < 0 || row >= len(b) || col >= len(b[0]) {
		return
	}

	if b[row][col].count > -1 {
		return
	}

	mineCount := 0
	for _, dr := range directions {
		newRow, newCol := row+dr.x, col+dr.y

		if 0 <= newRow && newRow < len(b) && 0 <= newCol && newCol < len(b[0]) && b[newRow][newCol].isMine {
			mineCount += 1
		}
	}

	if mineCount > 0 {
		b[row][col].count = mineCount
		return
	}

	b[row][col].count = 0

	for _, dr := range directions {
		revealCells(b, row+dr.x, col+dr.y)
	}
}

func (m model) Init() tea.Cmd {
	// Just returno `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	isGameFinished := m.state.win || m.state.lose

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
			if !isGameFinished && m.cursorPos.y > 0 {
				m.cursorPos.y--
			}

		// The "down" and "j" keys move the cursor down
		case "left", "l":
			if !isGameFinished && m.cursorPos.y < len(m.board[0])-1 {
				m.cursorPos.y++
			}

		case "down", "j":
			if !isGameFinished && m.cursorPos.x < len(m.board)-1 {
				m.cursorPos.x++
			}

		case "up", "k":
			if !isGameFinished && m.cursorPos.x > 0 {
				m.cursorPos.x--
			}

		case "f":
			if !isGameFinished {
				m.board[m.cursorPos.x][m.cursorPos.y].isFlag = !m.board[m.cursorPos.x][m.cursorPos.y].isFlag
			}

		case "enter":
			if !isGameFinished && m.board[m.cursorPos.x][m.cursorPos.y].isMine && !m.board[m.cursorPos.x][m.cursorPos.y].isFlag {
				m.state.lose = !m.state.lose
			} else if !isGameFinished {
				revealCells(m.board, m.cursorPos.x, m.cursorPos.y)
			}

			// if !isGameFinished && m.board[m.cursorPos.x][m.cursorPos.y].isMine && !m.board[m.cursorPos.x][m.cursorPos.y].isFlag {
			// 	m.state.lose = !m.state.lose
			// }

		case "r":
			if isGameFinished {
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
			} else if m.board[i][j].count > 0 {
				m.board[i][j].el = fmt.Sprint(m.board[i][j].count)
			} else if m.board[i][j].count == 0 {
				m.board[i][j].el = "."
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

	boardTable := ui.NewBoardTable(rows)

	s += fmt.Sprintf("%s\n", boardTable)
	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func StartGame() {
	p := tea.NewProgram(newModel(9, 9))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
