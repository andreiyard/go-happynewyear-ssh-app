package bubbletea

import "math/rand"

type Point struct {
	X, Y int // position
}

func (p Point) move(dx, dy int) Point {
	return Point{p.X + dx, p.Y + dy}
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
		switch {
		case pos.Y == maxY: // if at the bottom, place at the same pos
			newSnowflakes[pos] = char
		case s.exists(nextPoint): // if on top of next snowflake, try fall to the side
			// Check left, then right, as a default stay at the same pos
			if !s.exists(pos.move(-1, 0)) && !s.exists(pos.move(-1, 1)) && !s.exists(pos.move(-1, 2)) {
				newSnowflakes[pos.move(-1, 1)] = char
			} else if !s.exists(pos.move(1, 0)) && !s.exists(pos.move(1, 1)) && !s.exists(pos.move(1, 2)) {
				newSnowflakes[pos.move(1, 1)] = char
			} else {
				newSnowflakes[pos] = char
			}
		default: // If can fall, place at the next pos
			newSnowflakes[nextPoint] = char
		}
	}
	return newSnowflakes
}
