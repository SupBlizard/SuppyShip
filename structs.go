package main

import (
	"image/color"

	"github.com/faiface/pixel"
)

type player struct {
	pos  pixel.Vec
	vel  pixel.Vec
	dir  int16
	roll uint16
	acc  float64
	frc  float64

	// attrib
	power    uint8
	recharge bool
	alive    bool
	heat     uint16
	reload   uint16

	// obj
	hitbox circularHitbox
	sprite spriteSheet
	frag   fragInfo
	shield shipShield
}

type shipShield struct {
	active     bool
	prot       uint8
	protLength uint8
	sprite     spriteSheet
}

// Enemy Structure
type enemy struct {
	id        uint8
	pos       pixel.Vec
	vel       pixel.Vec
	rot       float64
	rotVel    float64
	health    uint16
	maxHealth uint16
	hitbox    circularHitbox
	sprite    spriteSheet
	frag      fragInfo
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
	cycleSpeed uint16
	scale      float64
	pos        [2]pixel.Rect
}

// A Fragment created upon destruction
type fragment struct {
	ID     [2]uint8
	pos    pixel.Vec
	vel    pixel.Vec
	rot    float64
	rotVel float64
	scale  float64
}

// Info about the object's fragments
type fragInfo struct {
	ID     uint8
	frags  uint8
	power  float64
	radius float64
	scale  float64
}

// Background star Structure
type star struct {
	pos   pixel.Vec
	scale float64
	phase int8
	shine int8
}

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

// Player Inputs Structure
type inputStruct struct {
	dir   pixel.Vec
	shoot bool
	roll  bool
}

// Circular hitbox Structure
type circularHitbox struct {
	radius float64
	offset pixel.Vec
}

type trailPart struct {
	pos  pixel.Vec
	mask color.RGBA
}
