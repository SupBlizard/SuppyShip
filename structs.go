package main

import "github.com/faiface/pixel"

// Basic physics struct
type physObj struct {
	pos pixel.Vec // position
	vel pixel.Vec // velocity
	acc float64   // acceleration
	frc float64   // friction
}

// Player Struct
type player struct {
	phys   physObj
	hitbox circularHitbox
	power  uint8
	sprite spriteSheet
}

// Player Inputs
type inputStruct struct {
	dir   pixel.Vec
	shoot bool
	roll  bool
}

// Enemy struct
type enemy struct {
	phys   physObj
	loaded bool
	hitbox circularHitbox
	health uint16
	sprite spriteSheet
	name   string
	id     int
}

// Circular hitbox
type circularHitbox struct {
	radius float64
	offset pixel.Vec
}

// Projectile structure
type projectile struct {
	id       uint8
	pos      pixel.Vec
	vel      pixel.Vec
	name     string
	loaded   bool
	friendly bool
	sprite   projectileSprite
}

// Projectile sprite
type projectileSprite struct {
	cycle      uint8
	cycleSpeed int
	scale      float64
	pos        [2]pixel.Rect
}
