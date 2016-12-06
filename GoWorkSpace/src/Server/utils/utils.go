package utils

import (
	"strconv"
)

type Point struct {
	X, Y    int
	H, G, F int
	Parent  *Point
}

func (p Point) String() string {
	return "[" + strconv.Itoa(p.X) + ", " + strconv.Itoa(p.Y) + ", " + strconv.Itoa(p.F) + "]"
}