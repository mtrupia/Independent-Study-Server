package main

import (
	"time"
)

func main() {
	var scene Scene
	scene.initScene(50, 50)
	scene.draw()
	scene.addWalls(10)
	scene.draw()
	initAstar(&scene)
	scene.draw()

	for {
		findPath(&scene)
		scene.draw()
		time.Sleep(50 * time.Millisecond)
	}
}
