package bubbletea

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

//TODO: Draw a christmas tree (snowflakes should cover it after some time)
//TODO: Stop generating snowflakes after some limit (or start deleting the old ones)

const fps = 6

// New snowflake amount per frame
const snowflakeRate = 2

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

type Point struct {
	x, y int // position
}

type Snowflakes map[Point]string

func (s Snowflakes) drawSnowflakes(cellbuffer *cellbuffer) {
	for pos, char := range s {
		cellbuffer.setChar(pos.x, pos.y, char)
	}
}

func (s Snowflakes) placeSnowflake(x, y int, char string) {
	s[Point{x, y}] = char
}

func (s Snowflakes) addSnowflake(x, y int, options []string) {
	char := options[rand.Intn(len(options))]
	s.placeSnowflake(x, y, char)
}

func (s Snowflakes) addRandomPosSnowflake(maxX int, options []string) {
	x := rand.Intn(maxX)
	s.addSnowflake(x, 0, options)
}

func (s Snowflakes) exists(p Point) bool {
	_, exists := s[p]
	return exists
}

func (s Snowflakes) advanceAll(maxY int) Snowflakes {
	newSnowflakes := make(Snowflakes)
	for pos, char := range s {
		nextPoint := Point{pos.x, pos.y + 1}
		if pos.y == maxY || s.exists(nextPoint) {
			newSnowflakes[pos] = char
		} else {
			newSnowflakes[nextPoint] = char
		}
	}
	return newSnowflakes
}

type model struct {
	cells          cellbuffer
	snowflakes     Snowflakes
	snowflakeChars []string
}

func NewModel(snowflakeChars []string) model {
	return model{snowflakeChars: snowflakeChars}
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
		m.snowflakes = make(Snowflakes)
		m.cells.init(msg.Width, msg.Height)
		return m, nil
	case tea.MouseMsg:
		//Add new snowflake at pos of mouse click
		m.snowflakes.addSnowflake(msg.X, msg.Y, m.snowflakeChars)
		return m, nil
	case frameMsg:
		if !m.cells.ready() {
			return m, nil
		}
		// Wipe last frame
		m.cells.wipe()

		// Create new snowflakes and advance existing
		maxX := m.cells.width() - 1
		maxY := m.cells.height() - 1
		for range snowflakeRate {
			m.snowflakes.addRandomPosSnowflake(maxX, m.snowflakeChars)
		}
		m.snowflakes = m.snowflakes.advanceAll(maxY)

		// Draw all snowflakes on the frame
		m.snowflakes.drawSnowflakes(&m.cells)
		return m, animate()
	default:
		return m, nil
	}
}

func (m model) View() string {
	return m.cells.String()
}
