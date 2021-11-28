package main

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/project-auxo/auxo/apollo/model/pendulum"
)

func main() {
	pixelgl.Run(pendulum.Run)
}
