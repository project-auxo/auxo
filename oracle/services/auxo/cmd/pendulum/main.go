package main

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/project-auxo/auxo/oracle/services/auxo/pendulum"
)

func main() {
	pixelgl.Run(pendulum.Run)
}
