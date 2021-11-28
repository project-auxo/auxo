package main

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/project-auxo/auxo/apollo/model/seek"
)

func main() {
	pixelgl.Run(seek.Run)
}
