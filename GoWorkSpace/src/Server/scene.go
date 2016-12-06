package main

import (
	"fmt"
)

type Scene struct {
	rows, cols int
	scene      [][]byte
}

func (s *Scene) initScene(p * Player) {
	s.rows = 20
	s.cols = 25

	s.scene = make([][]byte, s.rows)
	for i := 0; i < s.rows; i++ {
		s.scene[i] = make([]byte, s.cols)
		for j := 0; j < s.cols; j++ {
			if i == 0 || i == s.rows-1 || j == 0 || j == s.cols-1 {
				s.scene[i][j] = '#'
			} else {
				s.scene[i][j] = ' '
			}
			
			// add walls if exists
			if (i >= 1 && i <= s.rows-2) && (j >= 1 && j <= s.cols-2) {
				if (p.Field[i-1][j-1].Building.Id >= 20) {
					s.scene[i][j] = '#'
				} 
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