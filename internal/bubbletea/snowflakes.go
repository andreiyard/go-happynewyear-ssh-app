package bubbletea

import "math/rand"

type Point struct {
	X, Y int // position
}

type Snowflakes map[Point]string

func (s Snowflakes) drawSnowflakes(cellbuffer *cellbuffer) {
	for pos, char := range s {
		cellbuffer.setChar(pos.X, pos.Y, char)
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

func (s Snowflakes) deleteRandomN(n int) {
	count := 0
	for k := range s {
		if count >= n {
			return
		}
		delete(s, k)
		count++
	}
}

func (s Snowflakes) advanceAll(maxY int) Snowflakes {
	newSnowflakes := make(Snowflakes)
	for pos, char := range s {
		nextPoint := Point{pos.X, pos.Y + 1}
		if pos.Y == maxY || s.exists(nextPoint) {
			newSnowflakes[pos] = char
		} else {
			newSnowflakes[nextPoint] = char
		}
	}
	return newSnowflakes
}
