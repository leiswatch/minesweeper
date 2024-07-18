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

type cell struct {
	char      string
	isFlagged bool
	isMine    bool
	count     int
}

type board = [][]cell

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
	cursorChar string = "X"
	flagChar   string = "F"
	hiddenChar string = " "
	mineChar   string = "B"
	emptyChar  string = "."
)

var directions = [8]coordinates{
	{x: 1, y: 0},
	{x: -1, y: 0},
	{x: 0, y: 1},
	{x: 0, y: -1},
	{x: 1, y: -1},
	{x: 1, y: 1},
	{x: -1, y: 1},
	{x: -1, y: -1},
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
		board[i] = make([]cell, height)

		for j := range height {
			coords := coordinates{
				x: i,
				y: j,
			}

			board[i][j] = cell{
				isFlagged: false,
				isMine:    slices.Contains(bombPositions, coords),
				char:      hiddenChar,
				count:     -1,
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

func (m *model) revealCells(row, col int) {
	if row < 0 || col < 0 || row >= len(m.board) || col >= len(m.board[0]) {
		return
	}

	if m.board[row][col].count >= 0 {
		return
	}

	mineCount := 0
	for _, dr := range directions {
		newRow, newCol := row+dr.x, col+dr.y

		if newRow < 0 || newCol < 0 || newRow >= len(m.board) || newCol >= len(m.board[0]) {
			continue
		}

		if m.board[newRow][newCol].isMine {
			mineCount += 1
		}
	}

	m.board[row][col].count = mineCount

	if mineCount > 0 {
		return
	}

	for _, dr := range directions {
		m.revealCells(row+dr.x, col+dr.y)
	}
}

func (m *model) checkWin() {
	count := 0

	for i := range m.board {
		for j := range m.board[i] {
			if m.board[i][j].char == hiddenChar || m.board[i][j].isFlagged {
				count += 1
			}
		}
	}

	if count == 10 {
		m.state.win = true
	}
}

func (m model) Init() tea.Cmd {
	// Just returno `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	isGameFinished := m.state.win || m.state.lose
	isMine := m.board[m.cursorPos.x][m.cursorPos.y].isMine

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
				m.board[m.cursorPos.x][m.cursorPos.y].isFlagged = !m.board[m.cursorPos.x][m.cursorPos.y].isFlagged
				m.checkWin()
			}

		case "enter":
			if !isGameFinished && isMine {
				m.state.lose = !m.state.lose
			} else if !isGameFinished && !isMine {
				m.revealCells(m.cursorPos.x, m.cursorPos.y)
				m.checkWin()
			}

		case "r":
			if isGameFinished {
				m.board, m.state.lose, m.state.win = newBoard(9, 9), false, false
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := ""

	isGameFinished := m.state.win || m.state.lose

	for i := range len(m.board) {
		for j := range len(m.board[i]) {
			if m.state.lose && m.board[i][j].isMine {
				m.board[i][j].char = mineChar
			} else if m.cursorPos.x == i && m.cursorPos.y == j {
				m.board[i][j].char = cursorChar
			} else if m.board[i][j].isFlagged {
				m.board[i][j].char = flagChar
			} else if m.board[i][j].count > 0 {
				m.board[i][j].char = fmt.Sprint(m.board[i][j].count)
			} else if m.board[i][j].count == 0 {
				m.board[i][j].char = ui.RenderGrayText(emptyChar)
			} else {
				m.board[i][j].char = hiddenChar
			}
		}
	}

	rows := make([][]string, len(m.board))

	for i := range len(m.board) {
		rows[i] = make([]string, len(m.board[i]))

		for j := range len(m.board[i]) {
			rows[i][j] = m.board[i][j].char
		}
	}

	if m.state.win {
		s += ui.RenderGreenText("You won!")
	} else if m.state.lose {
		s += ui.RenderRedText("You lost!")
	} else {
		s += ui.RenderWhiteText("Minesweeper")
	}

	boardTable := ui.NewBoard(rows)

	s += fmt.Sprintf("\n%s\n", boardTable)
	if isGameFinished {
		s += fmt.Sprintf("\nPress %s to %s.", ui.RenderBoldText("r"), ui.RenderBlueText("reload"))
	}

	s += fmt.Sprintf("\nPress %s to %s.\n", ui.RenderBoldText("q"), ui.RenderRedText("quit"))

	return s
}

func StartGame() {
	p := tea.NewProgram(newModel(9, 9))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
