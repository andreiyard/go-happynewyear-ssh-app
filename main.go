package main

// A simple example demonstrating how to draw and animate on a cellular grid.
// Note that the cellbuffer implementation in this example does not support
// double-width runes.

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var snowflakeChars = [...]string{"*", "+", "."}

const fps = 60

type cellbuffer struct {
	cells  []string
	stride int
}

func (c *cellbuffer) init(w, h int) {
	if w == 0 {
		return
	}
	c.stride = w
	c.cells = make([]string, w*h)
	c.wipe()
}

func (c cellbuffer) setChar(x, y int, char string) {
	i := y*c.stride + x
	if i > len(c.cells)-1 || x < 0 || y < 0 || x >= c.width() || y >= c.height() {
		return
	}
	c.cells[i] = char
}

func (c *cellbuffer) wipe() {
	for i := range c.cells {
		c.cells[i] = " "
	}
}

func (c cellbuffer) width() int {
	return c.stride
}

func (c cellbuffer) height() int {
	h := len(c.cells) / c.stride
	if len(c.cells)%c.stride != 0 {
		h++
	}
	return h
}

func (c cellbuffer) ready() bool {
	return len(c.cells) > 0
}

func (c cellbuffer) String() string {
	var b strings.Builder
	for i := 0; i < len(c.cells); i++ {
		if i > 0 && i%c.stride == 0 && i < len(c.cells)-1 {
			b.WriteRune('\n')
		}
		b.WriteString(c.cells[i])
	}
	return b.String()
}

type frameMsg struct{}

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(_ time.Time) tea.Msg {
		return frameMsg{}
	})
}

type Snowflake struct {
	x, y int    // position
	char string // char from snowflakeChars
}

type model struct {
	cells      cellbuffer
	snowflakes []Snowflake
}

func (m model) Init() tea.Cmd {
	return animate()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		return m, nil
	case tea.WindowSizeMsg:
		// Reinit animation on window change
		m.snowflakes = make([]Snowflake, 0, 100)
		m.cells.init(msg.Width, msg.Height)
		return m, nil
	case tea.MouseMsg:
		// TODO: Implement snowflake creation at pos of mouse click
		return m, nil
	case frameMsg:
		if !m.cells.ready() {
			return m, nil
		}
		// Wipe last frame
		m.cells.wipe()

		// Create new snowflakes and advance existing
		m.cells.setChar(10, 11, "*")
		m.cells.setChar(11, 12, "*")
		m.cells.setChar(12, 14, "*")
		m.cells.setChar(15, 17, "*")

		// Render snowflakes on the frame
		//drawSnowflakes(&m.cells, &m.snowflakes)
		return m, animate()
	default:
		return m, nil
	}
}

func (m model) View() string {
	return m.cells.String()
}

func main() {
	m := model{}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
