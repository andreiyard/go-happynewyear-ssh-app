package bubbletea

import "testing"

func TestSnowflakes_Exists(t *testing.T) {
	tests := []struct {
		name       string
		snowflakes Snowflakes
		point      Point
		want       bool
	}{
		{
			name:       "exists at (0,0)",
			snowflakes: Snowflakes{{0, 0}: "*"},
			point:      Point{0, 0},
			want:       true,
		},
		{
			name:       "exists at (5,10)",
			snowflakes: Snowflakes{{5, 10}: "+"},
			point:      Point{5, 10},
			want:       true,
		},
		{
			name:       "does not exist",
			snowflakes: Snowflakes{{0, 0}: "*"},
			point:      Point{1, 1},
			want:       false,
		},
		{
			name:       "empty map",
			snowflakes: Snowflakes{},
			point:      Point{0, 0},
			want:       false,
		},
		{
			name: "multiple snowflakes",
			snowflakes: Snowflakes{
				{0, 0}: "*",
				{1, 1}: "+",
				{2, 2}: ".",
			},
			point: Point{1, 1},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.snowflakes.exists(tt.point); got != tt.want {
				t.Errorf("exists(%v) = %v, want %v", tt.point, got, tt.want)
			}
		})
	}
}

func TestSnowflakes_PlaceSnowflake(t *testing.T) {
	tests := []struct {
		name     string
		x        int
		y        int
		char     string
		wantChar string
	}{
		{
			name:     "place at origin",
			x:        0,
			y:        0,
			char:     "*",
			wantChar: "*",
		},
		{
			name:     "place at (10,20)",
			x:        10,
			y:        20,
			char:     "+",
			wantChar: "+",
		},
		{
			name:     "place dot",
			x:        5,
			y:        5,
			char:     ".",
			wantChar: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make(Snowflakes)
			s.placeSnowflake(tt.x, tt.y, tt.char)

			p := Point{tt.x, tt.y}
			if !s.exists(p) {
				t.Errorf("snowflake not placed at (%d,%d)", tt.x, tt.y)
			}
			if got := s[p]; got != tt.wantChar {
				t.Errorf("snowflake char = %q, want %q", got, tt.wantChar)
			}
		})
	}
}

func TestSnowflakes_AdvanceAll(t *testing.T) {
	tests := []struct {
		name       string
		snowflakes Snowflakes
		maxY       int
		want       Snowflakes
	}{
		{
			name:       "single snowflake falls",
			snowflakes: Snowflakes{{5, 3}: "*"},
			maxY:       10,
			want:       Snowflakes{{5, 4}: "*"},
		},
		{
			name:       "snowflake at bottom stays",
			snowflakes: Snowflakes{{5, 10}: "*"},
			maxY:       10,
			want:       Snowflakes{{5, 10}: "*"},
		},
		{
			name: "snowflake on top of another stays, bottom one moves",
			snowflakes: Snowflakes{
				{5, 3}: "*",
				{5, 4}: "+",
			},
			maxY: 10,
			want: Snowflakes{
				{5, 3}: "*",
				{5, 5}: "+",
			},
		},
		{
			name: "multiple snowflakes falling",
			snowflakes: Snowflakes{
				{0, 0}: "*",
				{5, 5}: "+",
				{9, 2}: ".",
			},
			maxY: 10,
			want: Snowflakes{
				{0, 1}: "*",
				{5, 6}: "+",
				{9, 3}: ".",
			},
		},
		{
			name: "mixed falling and stopped",
			snowflakes: Snowflakes{
				{3, 5}:  "*",
				{3, 6}:  "+",
				{3, 10}: ".",
				{7, 2}:  "*",
			},
			maxY: 10,
			want: Snowflakes{
				{3, 5}:  "*",
				{3, 7}:  "+",
				{3, 10}: ".",
				{7, 3}:  "*",
			},
		},
		{
			name:       "empty map",
			snowflakes: Snowflakes{},
			maxY:       10,
			want:       Snowflakes{},
		},
		{
			name: "stack building from bottom",
			snowflakes: Snowflakes{
				{5, 9}:  "*",
				{5, 10}: "+",
			},
			maxY: 10,
			want: Snowflakes{
				{5, 9}:  "*",
				{5, 10}: "+",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.snowflakes.advanceAll(tt.maxY)

			if len(got) != len(tt.want) {
				t.Errorf("advanceAll() returned %d snowflakes, want %d", len(got), len(tt.want))
			}

			for point, wantChar := range tt.want {
				if !got.exists(point) {
					t.Errorf("expected snowflake at %v not found", point)
					continue
				}
				if gotChar := got[point]; gotChar != wantChar {
					t.Errorf("snowflake at %v = %q, want %q", point, gotChar, wantChar)
				}
			}

			for point := range got {
				if _, exists := tt.want[point]; !exists {
					t.Errorf("unexpected snowflake at %v", point)
				}
			}
		})
	}
}

func TestSnowflakes_DeleteRandomN(t *testing.T) {
	tests := []struct {
		name        string
		initial     Snowflakes
		n           int
		wantRemains int
	}{
		{
			name: "delete fewer than exist",
			initial: Snowflakes{
				{0, 0}: "*",
				{1, 1}: "+",
				{2, 2}: ".",
				{3, 3}: "*",
			},
			n:           2,
			wantRemains: 2,
		},
		{
			name: "delete more than exist",
			initial: Snowflakes{
				{0, 0}: "*",
				{1, 1}: "+",
			},
			n:           5,
			wantRemains: 0,
		},
		{
			name: "delete exactly all",
			initial: Snowflakes{
				{0, 0}: "*",
				{1, 1}: "+",
				{2, 2}: ".",
			},
			n:           3,
			wantRemains: 0,
		},
		{
			name:        "delete from empty",
			initial:     Snowflakes{},
			n:           5,
			wantRemains: 0,
		},
		{
			name: "delete zero",
			initial: Snowflakes{
				{0, 0}: "*",
				{1, 1}: "+",
			},
			n:           0,
			wantRemains: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make(Snowflakes)
			for k, v := range tt.initial {
				s[k] = v
			}

			s.deleteRandomN(tt.n)

			if got := len(s); got != tt.wantRemains {
				t.Errorf("after deleteRandomN(%d), len = %d, want %d", tt.n, got, tt.wantRemains)
			}
		})
	}
}

func TestSnowflakes_DrawSnowflakes(t *testing.T) {
	tests := []struct {
		name       string
		snowflakes Snowflakes
		width      int
		height     int
		wantAt     []struct{ x, y int; char string }
	}{
		{
			name: "single snowflake",
			snowflakes: Snowflakes{
				{2, 1}: "*",
			},
			width:  5,
			height: 3,
			wantAt: []struct{ x, y int; char string }{
				{2, 1, "*"},
			},
		},
		{
			name: "multiple snowflakes",
			snowflakes: Snowflakes{
				{0, 0}: "*",
				{4, 2}: "+",
				{2, 1}: ".",
			},
			width:  5,
			height: 3,
			wantAt: []struct{ x, y int; char string }{
				{0, 0, "*"},
				{4, 2, "+"},
				{2, 1, "."},
			},
		},
		{
			name:       "empty snowflakes",
			snowflakes: Snowflakes{},
			width:      5,
			height:     3,
			wantAt:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c cellbuffer
			c.init(tt.width, tt.height)

			tt.snowflakes.drawSnowflakes(&c)

			for _, want := range tt.wantAt {
				i := want.y*c.stride + want.x
				if got := c.cells[i]; got != want.char {
					t.Errorf("cell at (%d,%d) = %q, want %q", want.x, want.y, got, want.char)
				}
			}
		})
	}
}
