package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

const (
	WINX float64 = 512
	WINY float64 = 768

	BOUNDARY_STRENGTH       float64 = 2
	AXIS_DEADZONE           float64 = 0.1
	DEFAULT_GLOBAL_VELOCITY float64 = 10

	ROLL_COOLDOWN uint16 = 20
	ONYX_COOLDOWN uint16 = 60

	ONYX_CLUSTER_REQUIREMENT uint16  = 7
	ONYX_CLUSTER_RADIUS      float64 = 30

	// Allocation sizes
	ENEMY_ALLOC_SIZE  uint16 = 16
	BULLET_ALLOC_SIZE uint16 = 256

	// Stars
	STAR_MAX_PHASE   int8    = 5
	STAR_PHASES      int8    = (STAR_MAX_PHASE) * 2
	STAR_DISTANCE    float64 = 65
	STAR_RANDOMNESS  int32   = 60
	STAR_SIZE        float64 = 5
	STARFIELD_NUMBER uint8   = 2
)

var (
	currentLevel   uint8
	globalVelocity float64    = DEFAULT_GLOBAL_VELOCITY
	globalAcc      [2]float64 = [2]float64{1.4, 0.6}

	// Input
	input       = inputStruct{}
	inputLookup = [4]pixel.Vec{
		pixel.V(0, 1),
		pixel.V(-1, 0),
		pixel.V(0, -1),
		pixel.V(1, 0),
	}

	gunCooldown uint16
	reloadDelay uint32 = 4

	// Border values (top, bottom, sides)
	windowBorder = [3]float64{0, 0, 0}
	forceBorder  = [3]float64{400, 150, 70}
	spawnBorder  = [3]float64{-300, -50, -100}

	// Allocation
	enemies     [ENEMY_ALLOC_SIZE]enemy
	projectiles [BULLET_ALLOC_SIZE]projectile
	shipTrail   []trailPart = make([]trailPart, 0, 128)

	// Loaded objects
	loadedEnemies     []uint16 = make([]uint16, 0, ENEMY_ALLOC_SIZE)
	loadedProjectiles []uint16 = make([]uint16, 0, BULLET_ALLOC_SIZE)

	// Spritesheets
	projectileSheet  pixel.Picture = loadPicture("assets/projectile-spritesheet.png")
	trailSpritesheet pixel.Picture = loadPicture("assets/trail.png")
	starSheet        pixel.Picture = loadPicture("assets/star-spritesheet.png")

	// Batches
	trailBatch      *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, trailSpritesheet)
	projectileBatch *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, projectileSheet)

	// Sprites
	starSprites [STAR_MAX_PHASE + 1]*pixel.Sprite
)

// Calculate how far into the border something is
func findBorderDepth(pos float64, borderRange float64) float64 { return 1 - pos/borderRange }

// Return 1 or -1 for the signbit
func signbit(x float64) float64 { return x / math.Abs(x) }

// Return true when frameCount is a multiple of x
func skipFrames(x uint32) bool { return frameCount%x == 0 }

// Return a random vector
func randomVector(limit int32) pixel.Vec {
	return pixel.V(float64(rand.Int31()%limit), float64(rand.Int31()%limit))
}

// Check if pos is in bounds
func inBounds(pos pixel.Vec, boundaryRange [3]float64) pixel.Vec {
	var boundCollision pixel.Vec = pixel.ZV
	if pos.Y >= WINY-boundaryRange[0] {
		boundCollision.Y = -1
	} else if pos.Y <= boundaryRange[1] {
		boundCollision.Y = 1
	}
	if pos.X >= WINX-boundaryRange[2] {
		boundCollision.X = -1
	} else if pos.X <= boundaryRange[2] {
		boundCollision.X = 1
	}

	return boundCollision
}

// Return all of the projectiles in a certain radius
func projectilesInRadius(point pixel.Vec, radius float64, friendliness bool) ([]uint16, uint16) {
	var inside []uint16
	var count uint16 = 0

	// Loop through loaded indexes
	for _, i := range loadedProjectiles {
		if projectiles[i].friendly == friendliness && projectiles[i].pos.Sub(point).Len() < radius {
			inside = append(inside, i)
			count++
		}
	}

	return inside, count
}
