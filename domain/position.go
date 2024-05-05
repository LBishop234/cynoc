package domain

import "github.com/davecgh/go-spew/spew"

type Position struct {
	x int
	y int
}

func NewPosition(x, y int) Position {
	return Position{
		x: x,
		y: y,
	}
}

func (p Position) X() int {
	return p.x
}

func (p Position) Y() int {
	return p.y
}

func (p Position) Prettify() string {
	return spew.Sdump(p)
}
