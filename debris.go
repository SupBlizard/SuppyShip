package main

import "github.com/faiface/pixel"

var debrisTypes = [3]debris{
	{
		pos:    pixel.Vec{},
		vel:    pixel.Vec{},
		rot:    0,
		rotVel: 0,
		scale:  0,
		sprite: [3]pixel.Rect{},
	},
}
