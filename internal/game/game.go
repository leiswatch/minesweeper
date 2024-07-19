package game

import (
	"fmt"
	"math/rand/v2"
	"os"
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

func newBoard(width, height, numMines int) board {
	board := make(board, width)
	positions := make([]coordinates, 0, width*height)
	seed := uint64(time.Now().UnixNano())
	r := rand.New(rand.NewPCG(seed, seed))

	for i := range width {
		board[i] = make([]cell, height)

		for j := range height {
			coords := coordinates{
				x: i,
				y: j,
			}
			positions = append(positions, coords)

			board[i][j] = cell{
				isFlagged: false,
				isMine:    false,
				char:      hiddenChar,
				count:     -1,
			}
		}
	}

	r.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	for i := 0; i < numMines; i++ {
		board[positions[i].x][positions[i].y].isMine = true
	}

	return board
}

func newModel(width, height, numMines int) model {
	return model{
		board: newBoard(width, height, numMines),
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
	cells := make([]coordinates, 0)

	for i := range m.board {
		for j := range m.board[i] {
			if m.board[i][j].char == hiddenChar || m.board[i][j].isFlagged {
				cells = append(cells, coordinates{
					x: i,
					y: j,
				})
			}
		}
	}

	if len(cells) == 10 {
		for i := range cells {
			if !m.board[cells[i].x][cells[i].y].isMine {
				return
			}
		}

		m.state.win = true
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currX := m.cursorPos.x
	currY := m.cursorPos.y
	isGameFinished := m.state.win || m.state.lose
	isMine := m.board[currX][currY].isMine

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "h":
			if !isGameFinished && currY > 0 {
				m.cursorPos.y--
			}
		case "left", "l":
			if !isGameFinished && currY < len(m.board[0])-1 {
				m.cursorPos.y++
			}
		case "down", "j":
			if !isGameFinished && currX < len(m.board)-1 {
				m.cursorPos.x++
			}
		case "up", "k":
			if !isGameFinished && currX > 0 {
				m.cursorPos.x--
			}
		case "f":
			if !isGameFinished {
				m.board[currX][currY].isFlagged = !m.board[currX][currY].isFlagged
				m.checkWin()
			}
		case "enter":
			if !isGameFinished && isMine {
				m.state.lose = !m.state.lose
			} else if !isGameFinished && !isMine {
				m.revealCells(currX, currY)
				m.checkWin()
			}
		case "r":
			if isGameFinished {
				m.board, m.state.lose, m.state.win = newBoard(9, 9, 10), false, false
			}
		}
	}

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
				m.board[i][j].char = ui.RenderText(flagChar, ui.PinkColor, false)
			} else if m.board[i][j].count > 0 {
				m.board[i][j].char = fmt.Sprint(m.board[i][j].count)
			} else if m.board[i][j].count == 0 {
				m.board[i][j].char = ui.RenderText(emptyChar, ui.GrayColor, false)
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
		s += ui.RenderText("You won!", ui.GreenColor, true)
	} else if m.state.lose {
		s += ui.RenderText("You lost!", ui.RedColor, true)
	} else {
		s += ui.RenderText("Minesweeper", ui.WhiteColor, true)
	}

	boardTable := ui.NewBoard(rows)

	s += fmt.Sprintf("\n%s\n", boardTable)

	if isGameFinished {
		s += fmt.Sprintf("\nPress %s to %s.", ui.RenderText("r", ui.BlueColor, true), ui.RenderText("reload", ui.WhiteColor, true))
	}

	s += fmt.Sprintf("\nPress %s to %s.\n", ui.RenderText("q", ui.RedColor, true), ui.RenderText("quit", ui.WhiteColor, true))

	return s
}

func StartGame() {
	p := tea.NewProgram(newModel(9, 9, 10))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
