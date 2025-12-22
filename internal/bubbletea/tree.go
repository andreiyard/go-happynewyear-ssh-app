package bubbletea

type Tree map[Point]string

func newTree(w, h, windowW, windowH int) Tree {
	tree := make(Tree)

	// Our tree must start at the bottom middle
	yOffset := windowH - h
	xOffset := windowW/2 - w/2

	// Separate leafy parts and trunk
	trunkHeight := h / 6
	trunkStart := h - trunkHeight

	// TODO: Improve the design (it is quite simple at the moment)
	// TODO: Draw a star at the top of a tree
	// Also maybe some lights (or animate them)

	// Top part
	for j := range trunkStart {
		// Should equal 1 at the top, and "w" at the bottom
		widthOfCurrentLayer := 1 + (j-1)*(w-1)/(trunkStart-1)
		// Add offset to center the layer
		offset := (w - widthOfCurrentLayer) / 2
		for i := range widthOfCurrentLayer {
			tree[Point{i + xOffset + offset, j + yOffset}] = "▲"
		}
	}

	// Bottom part
	spacing := w / 3
	trunkLeft := spacing
	trunkRight := w - spacing
	for j := trunkStart; j < h; j++ {
		for i := trunkLeft; i < trunkRight; i++ {
			tree[Point{i + xOffset, j + yOffset}] = "|"
		}
	}

	return tree
}

func (t Tree) draw(cellbuffer *cellbuffer) {
	for p, char := range t {
		cellbuffer.setChar(p.X, p.Y, char)
	}
}
