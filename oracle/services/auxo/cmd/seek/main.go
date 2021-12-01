package main

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/project-auxo/auxo/oracle/services/auxo/seek"
)

func main() {
	pixelgl.Run(seek.Run)
}
