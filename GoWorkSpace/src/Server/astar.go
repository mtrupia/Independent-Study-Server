package main

import (
	"Server/utils"
	"math"
)

var origin, dest, path utils.Point
var openList, closeList []utils.Point

// Set the origin point
func setOrig(s *Scene) {
	origin = utils.Point{18, 12, 0, 0, 0, nil}
	s.scene[origin.X][origin.Y] = 'A'
}

// Set the destination point '
func setDest(s *Scene) {
	dest = utils.Point{1, 12, 0, 0, 0, nil}
	s.scene[dest.X][dest.Y] = 'B'
}

// Init origin, destination. Put the origin point into the openlist by the way
func initAstar(s *Scene) {
	openList = nil
	closeList = nil
	first = true

	setOrig(s)
	setDest(s)
	openList = append(openList, origin)
}

func findPath(s *Scene) utils.Point{
	current := getFMin()
	addToCloseList(current, s)
	walkable := getWalkable(current, s)
	for _, p := range walkable {
		addToOpenList(p)
	}
	
	if  hasPath() {
		findPath(s)
	}
	return path
}

func hasPath() bool {
	if len(openList) == 0 {
		return false
	}
	return true
}

func getFMin() utils.Point {
	index := 0
	for i, p := range openList {
		if (i > 0) && (p.F <= openList[index].F) {
			index = i
		}
	}
	return openList[index]
}

func getWalkable(p utils.Point, s *Scene) []utils.Point {
	var around []utils.Point
	row, col := p.X, p.Y
	left := s.scene[row][col-1]
	up := s.scene[row-1][col]
	right := s.scene[row][col+1]
	down := s.scene[row+1][col]
	if (left == ' ') || (left == 'B') {
		around = append(around, utils.Point{row, col - 1, 0, 0, 0, &p})
	}
	if (up == ' ') || (up == 'B') {
		around = append(around, utils.Point{row - 1, col, 0, 0, 0, &p})
	}
	if (right == ' ') || (right == 'B') {
		around = append(around, utils.Point{row, col + 1, 0, 0, 0, &p})
	}
	if (down == ' ') || (down == 'B') {
		around = append(around, utils.Point{row + 1, col, 0, 0, 0, &p})
	}
	return around
}

func addToOpenList(p utils.Point) {
	updateWeight(&p)
	if checkExist(p, closeList) {
		return
	}
	if !checkExist(p, openList) {
		openList = append(openList, p)
	} else {
		if openList[findPoint(p, openList)].F > p.F { //New path found
			openList[findPoint(p, openList)].Parent = p.Parent
		}
	}
}

// Update G, H, F of the point
func updateWeight(p *utils.Point) {
	if checkRelativePos(*p) == 1 {
		p.G = p.Parent.G + 10
	} else {
		p.G = p.Parent.G + 14
	}
	absx := (int)(math.Abs((float64)(dest.X - p.X)))
	absy := (int)(math.Abs((float64)(dest.Y - p.Y)))
	p.H = (absx + absy) * 10
	p.F = p.G + p.H
}

func removeFromOpenList(p utils.Point) {
	index := findPoint(p, openList)
	if index == -1 {
	}
	openList = append(openList[:index], openList[index+1:]...)
}

func addToCloseList(p utils.Point, s *Scene){
	removeFromOpenList(p)
	if (p.X == dest.X) && (p.Y == dest.Y) {
		generatePath(p, s)
		//s.draw()
	}
	// if (p.Parent != nil) && (checkRelativePos(p) == 2) {
	// 	parent := p.Parent
	// 	//rdblck := s.scene[p.X][parent.Y] | s.scene[parent.X][p.Y]
	// 	//fmt.Printf("%c\n", rdblck)
	// 	if (s.scene[p.X][parent.Y] == '#') || (s.scene[parent.X][p.Y] == '#') {
	// 		return
	// 	}
	// }
	if s.scene[p.X][p.Y] != 'A' {
		s.scene[p.X][p.Y] = '·'
	}
	closeList = append(closeList, p)
}

func checkExist(p utils.Point, arr []utils.Point) bool {
	for _, point := range arr {
		if p.X == point.X && p.Y == point.Y {
			return true
		}
	}
	return false
}

func findPoint(p utils.Point, arr []utils.Point) int {
	for index, point := range arr {
		if p.X == point.X && p.Y == point.Y {
			return index
		}
	}

	return -1
}

func checkRelativePos(p utils.Point) int {
	parent := p.Parent
	hor := (int)(math.Abs((float64)(p.X - parent.X)))
	ver := (int)(math.Abs((float64)(p.Y - parent.Y)))
	return hor + ver
}

var first bool = true
func generatePath(p utils.Point, s *Scene) {
	if first {
		path = p
		first = false
	}
	if (s.scene[p.X][p.Y] != 'A') && (s.scene[p.X][p.Y] != 'B') {
		s.scene[p.X][p.Y] = '*'
	}
	if p.Parent != nil {
		generatePath(*(p.Parent), s)
	}
}
