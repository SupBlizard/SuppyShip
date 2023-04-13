package main

import (
	"math"

	"github.com/faiface/pixel"
)

// Globals
const (
	BOUNDARY_STRENGTH       float64 = 2
	AXIS_DEADZONE           float64 = 0.1
	DEFAULT_GLOBAL_VELOCITY float64 = 10
	ROLL_COOLDOWN           uint16  = 20
	ENEMY_ALLOC_SIZE        uint8   = 16
)

var (
	frameCount   int
	currentLevel uint8
	winsize      pixel.Vec = pixel.V(512, 768)

	// Top Bottom Sides
	windowBorder = [3]float64{0, 0, 0}
	forceBorder  = [3]float64{400, 150, 70}
	spawnBorder  = [3]float64{-300, -50, -100}

	globalAcc      [2]float64 = [2]float64{1.4, 0.6}
	globalVelocity float64    = DEFAULT_GLOBAL_VELOCITY

	rollCooldown uint16
	gunCooldown  int
	reloadDelay  int = 4

	enemies       [ENEMY_ALLOC_SIZE]enemy
	loadedEnemies []uint8 = make([]uint8, 0, ENEMY_ALLOC_SIZE)

	input       = inputStruct{}
	inputLookup = [4]pixel.Vec{
		pixel.V(0, 1),
		pixel.V(-1, 0),
		pixel.V(0, -1),
		pixel.V(1, 0),
	}
)

// Calculate how far into the border something is
func findBorderDepth(pos float64, borderRange float64) float64 { return 1 - pos/borderRange }

// Return 1 or -1 for the signbit
func signbit(x float64) float64 { return x / math.Abs(x) }

// Return true when frameCount is a multiple of x
func skipFrames(x int) bool { return frameCount%x == 0 }

// Check if pos is in bounds
func inBounds(pos pixel.Vec, boundaryRange [3]float64) pixel.Vec {
	var boundCollision pixel.Vec = pixel.ZV
	if pos.Y >= winsize.Y-boundaryRange[0] {
		boundCollision.Y = -1
	} else if pos.Y <= boundaryRange[1] {
		boundCollision.Y = 1
	}
	if pos.X >= winsize.X-boundaryRange[2] {
		boundCollision.X = -1
	} else if pos.X <= boundaryRange[2] {
		boundCollision.X = 1
	}

	return boundCollision
}
