package main

import "github.com/faiface/pixel"

// Basic physics values
type physObj struct {
	pos pixel.Vec // position
	vel pixel.Vec // velocity
	acc float64   // acceleration
	frc float64   // friction
}

// Player values
type player struct {
	phys   physObj
	hitbox circularHitbox
	power  uint8
	sprite spriteSheet
}

// Generic enemy struct
type enemy struct {
	phys   physObj
	loaded bool
	hitbox circularHitbox
	health uint16
	sprite spriteSheet
	name   string
	id     int
}

// Player Inputs
type inputStruct struct {
	dir   pixel.Vec
	shoot bool
	roll  bool
}

// Circular hitbox
type circularHitbox struct {
	radius float64
	offset pixel.Vec
}
