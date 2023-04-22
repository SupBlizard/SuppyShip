package main

import (
	"image/color"

	"github.com/faiface/pixel"
)

// Spritesheet Structure
type spriteSheet struct {
	offset      uint16
	cycle       uint16
	cycleNumber uint16
	cycleSpeed  uint16
	current     uint16
	scale       float64
	sheet       []*pixel.Sprite
}

// Player Structure
type player struct {
	pos    pixel.Vec
	vel    pixel.Vec
	acc    float64
	frc    float64
	hitbox circularHitbox
	power  uint8
	sprite spriteSheet
}

// Player Inputs Structure
type inputStruct struct {
	dir   pixel.Vec
	shoot bool
	roll  bool
}

// Enemy Structure
type enemy struct {
	pos       pixel.Vec
	vel       pixel.Vec
	acc       float64
	frc       float64
	hitbox    circularHitbox
	health    uint16
	maxHealth uint16
	sprite    spriteSheet
	name      string
	id        uint8
}

// Circular hitbox Structure
type circularHitbox struct {
	radius float64
	offset pixel.Vec
}

// Projectile Structure
type projectile struct {
	id       uint8
	pos      pixel.Vec
	vel      pixel.Vec
	dmg      uint16
	friendly bool
	sprite   projectileSprite
}

// Projectile sprite Structure
type projectileSprite struct {
	cycle      uint8
	cycleSpeed uint32
	scale      float64
	pos        [2]pixel.Rect
}

// Background star Structure
type star struct {
	pos   pixel.Vec
	phase int8
	shine int8
}

type trailPart struct {
	pos  pixel.Vec
	mask color.RGBA
}

type debris struct {
	pos    pixel.Vec
	vel    pixel.Vec
	rot    float64
	rotVel float64
	scale  float64
	sprite [3]pixel.Rect
}
