package bubbletea

import (
	"strings"
	"testing"
)

func TestCellbuffer_Init(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		wantW  int
		wantH  int
		wantN  int
	}{
		{
			name:   "normal dimensions",
			width:  10,
			height: 5,
			wantW:  10,
			wantH:  5,
			wantN:  50,
		},
		{
			name:   "square",
			width:  5,
			height: 5,
			wantW:  5,
			wantH:  5,
			wantN:  25,
		},
		{
			name:   "wide",
			width:  20,
			height: 2,
			wantW:  20,
			wantH:  2,
			wantN:  40,
		},
		{
			name:   "tall",
			width:  3,
			height: 10,
			wantW:  3,
			wantH:  10,
			wantN:  30,
		},
		{
			name:   "zero width",
			width:  0,
			height: 5,
			wantW:  0,
			wantH:  0,
			wantN:  0,
		},
		{
			name:   "zero height",
			width:  5,
			height: 0,
			wantW:  5,
			wantH:  0,
			wantN:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c cellbuffer
			c.init(tt.width, tt.height)

			if got := c.width(); got != tt.wantW {
				t.Errorf("width() = %v, want %v", got, tt.wantW)
			}
			if tt.width > 0 {
				if got := c.height(); got != tt.wantH {
					t.Errorf("height() = %v, want %v", got, tt.wantH)
				}
			}
			if got := len(c.cells); got != tt.wantN {
				t.Errorf("len(cells) = %v, want %v", got, tt.wantN)
			}

			if tt.wantN > 0 {
				for i, cell := range c.cells {
					if cell != " " {
						t.Errorf("cell[%d] = %q, want \" \"", i, cell)
					}
				}
			}
		})
	}
}

func TestCellbuffer_SetChar(t *testing.T) {
	tests := []struct {
		name     string
		x        int
		y        int
		char     string
		wantSet  bool
		wantChar string
	}{
		{
			name:     "valid position (0,0)",
			x:        0,
			y:        0,
			char:     "*",
			wantSet:  true,
			wantChar: "*",
		},
		{
			name:     "valid position middle",
			x:        5,
			y:        2,
			char:     "+",
			wantSet:  true,
			wantChar: "+",
		},
		{
			name:     "valid position bottom-right",
			x:        9,
			y:        4,
			char:     ".",
			wantSet:  true,
			wantChar: ".",
		},
		{
			name:     "out of bounds x negative",
			x:        -1,
			y:        0,
			char:     "*",
			wantSet:  false,
			wantChar: " ",
		},
		{
			name:     "out of bounds y negative",
			x:        0,
			y:        -1,
			char:     "*",
			wantSet:  false,
			wantChar: " ",
		},
		{
			name:     "out of bounds x too large",
			x:        10,
			y:        0,
			char:     "*",
			wantSet:  false,
			wantChar: " ",
		},
		{
			name:     "out of bounds y too large",
			x:        0,
			y:        5,
			char:     "*",
			wantSet:  false,
			wantChar: " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c cellbuffer
			c.init(10, 5)
			c.setChar(tt.x, tt.y, tt.char)

			if tt.wantSet {
				i := tt.y*c.stride + tt.x
				if c.cells[i] != tt.wantChar {
					t.Errorf("cells[%d] = %q, want %q", i, c.cells[i], tt.wantChar)
				}
			}
		})
	}
}

func TestCellbuffer_Wipe(t *testing.T) {
	var c cellbuffer
	c.init(5, 3)

	c.setChar(0, 0, "*")
	c.setChar(2, 1, "+")
	c.setChar(4, 2, ".")

	c.wipe()

	for i, cell := range c.cells {
		if cell != " " {
			t.Errorf("cell[%d] = %q, want \" \" after wipe", i, cell)
		}
	}
}

func TestCellbuffer_String(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		chars  []struct {
			x, y int
			char string
		}
		want string
	}{
		{
			name:   "3x2 grid with stars",
			width:  3,
			height: 2,
			chars: []struct {
				x, y int
				char string
			}{
				{0, 0, "*"},
				{2, 0, "*"},
				{1, 1, "+"},
			},
			want: "* *\n + ",
		},
		{
			name:   "2x2 grid empty",
			width:  2,
			height: 2,
			chars:  nil,
			want:   "  \n  ",
		},
		{
			name:   "single row",
			width:  5,
			height: 1,
			chars: []struct {
				x, y int
				char string
			}{
				{0, 0, "a"},
				{4, 0, "b"},
			},
			want: "a   b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c cellbuffer
			c.init(tt.width, tt.height)

			for _, ch := range tt.chars {
				c.setChar(ch.x, ch.y, ch.char)
			}

			got := c.String()
			if got != tt.want {
				t.Errorf("String() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestCellbuffer_Ready(t *testing.T) {
	tests := []struct {
		name      string
		init      bool
		width     int
		height    int
		wantReady bool
	}{
		{
			name:      "initialized buffer",
			init:      true,
			width:     5,
			height:    5,
			wantReady: true,
		},
		{
			name:      "uninitialized buffer",
			init:      false,
			wantReady: false,
		},
		{
			name:      "zero width init",
			init:      true,
			width:     0,
			height:    5,
			wantReady: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c cellbuffer
			if tt.init {
				c.init(tt.width, tt.height)
			}

			if got := c.ready(); got != tt.wantReady {
				t.Errorf("ready() = %v, want %v", got, tt.wantReady)
			}
		})
	}
}

func TestCellbuffer_Dimensions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "10x5",
			width:  10,
			height: 5,
		},
		{
			name:   "1x1",
			width:  1,
			height: 1,
		},
		{
			name:   "100x50",
			width:  100,
			height: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c cellbuffer
			c.init(tt.width, tt.height)

			if got := c.width(); got != tt.width {
				t.Errorf("width() = %v, want %v", got, tt.width)
			}
			if got := c.height(); got != tt.height {
				t.Errorf("height() = %v, want %v", got, tt.height)
			}

			lines := strings.Split(c.String(), "\n")
			if got := len(lines); got != tt.height {
				t.Errorf("String() has %v lines, want %v", got, tt.height)
			}
			if tt.height > 0 {
				if got := len(lines[0]); got != tt.width {
					t.Errorf("String() first line has %v chars, want %v", got, tt.width)
				}
			}
		})
	}
}
