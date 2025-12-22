package bubbletea

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type frameMsg struct{}

func (m model) animate() tea.Cmd {
	return tea.Tick(time.Second/time.Duration(m.fps), func(_ time.Time) tea.Msg {
		return frameMsg{}
	})
}

type model struct {
	fps            int
	snowflakeRate  int
	snowflakeLimit int
	cells          cellbuffer
	snowflakes     Snowflakes
	snowflakeChars []string
	tree           Tree
}

func NewModel(fps, snowflakeRate, snowflakeLimit int, snowflakeChars []string) model {
	return model{
		fps:            fps,
		snowflakeRate:  snowflakeRate,
		snowflakeLimit: snowflakeLimit,
		snowflakeChars: snowflakeChars,
	}
}

func (m *model) reinit(w, h int) {
	treeHeight := h - h/4
	treeWidth := treeHeight
	m.tree = newTree(treeWidth, treeHeight, w, h)
	m.snowflakes = make(Snowflakes)
	m.cells.init(w, h)
}

func (m model) Init() tea.Cmd {
	return m.animate()
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
		// Reinit animation and tree on window change
		m.reinit(msg.Width, msg.Height)
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
		for range m.snowflakeRate {
			m.snowflakes.addRandomPosSnowflake(maxX, m.snowflakeChars)
		}
		m.snowflakes = m.snowflakes.advanceAll(maxY)

		// Delete random snowflakes if above limit (limit=0 means no limit)
		if m.snowflakeLimit != 0 {
			excessSnowflakes := len(m.snowflakes) - m.snowflakeLimit
			if excessSnowflakes > 0 {
				m.snowflakes.deleteRandomN(excessSnowflakes)
			}
		}
		// Draw all snowflakes on the frame
		m.snowflakes.drawSnowflakes(&m.cells)

		// Draw a tree on the frame
		m.tree.draw(&m.cells)

		return m, m.animate()
	default:
		return m, nil
	}
}

func (m model) View() string {
	return m.cells.String()
}
