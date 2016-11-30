package main

import (
	"fmt"
	"github.com/dyzdyz010/Golang-AStar/utils"
	"strconv"
)

type Scene struct {
	rows, cols int
	scene      [][]byte
}

func (s *Scene) initScene(rows int, cols int) {
	s.rows = rows
	s.cols = cols

	s.scene = make([][]byte, s.rows)
	for i := 0; i < s.rows; i++ {
		s.scene[i] = make([]byte, s.cols)
		for j := 0; j < s.cols; j++ {
			if i == 0 || i == s.rows-1 || j == 0 || j == s.cols-1 {
				s.scene[i][j] = '#'
			} else {
				s.scene[i][j] = ' '
			}
		}
	}
}

func (s *Scene) draw() {
	var color string
	var path string
	for i := 0; i < s.rows; i++ {
		for j := 0; j < s.cols; j++ {
			switch s.scene[i][j] {
			case '#':
				color += "# "		// wall
			case 'A':
				color += "A "		// Origin
			case 'B':
				color += "B "		// Destination
			case '*':
				color += "* "		// Path
				path += "["
				path += strconv.Itoa(i)
				path += ":"
				path += strconv.Itoa(j)
				path += "] "
			case ' ':
				color += "· "		// valid path
			default:
				color += "? "		// Finding Path
			}
		}
		color += "\n"
	}
	fmt.Printf("%s\n", color)
	fmt.Printf("%s\n", path)
}

func (s *Scene) addWalls(num int) {
	for i := 0; i < num; i++ {
		ori := utils.GetRandInt(2)
		length := utils.GetRandInt(16) + 1
		row := utils.GetRandInt(s.rows)
		col := utils.GetRandInt(s.cols)
		switch ori {
		case 0:
			for i := 0; i < length; i++ {
				if col+i >= s.cols {
					break
				}
				s.scene[row][col+i] = '#'
			}

		case 1:
			for i := 0; i < length; i++ {
				if row+i >= s.rows {
					break
				}
				s.scene[row+i][col] = '#'
			}
		}
	}
}
