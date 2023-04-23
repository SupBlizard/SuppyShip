package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

const (
	TITLE string = "Suppy Ship"

	WINX float64 = 512
	WINY float64 = 768

	REVOLUTION float64 = math.Pi * 2

	BOUNDARY_STRENGTH       float64 = 2
	AXIS_DEADZONE           float64 = 0.1
	DEFAULT_GLOBAL_VELOCITY float64 = 10

	ROLL_COOLDOWN            uint16  = 20
	ONYX_COOLDOWN            uint16  = 60
	ONYX_CLUSTER_REQUIREMENT uint16  = 7
	ONYX_CLUSTER_RADIUS      float64 = 30

	PROJ_SPRITES   uint8 = 2
	DEBRIS_SPRITES uint8 = 3

	// Allocation sizes
	ENEMY_ALLOC_SIZE  uint16 = 16
	DEBRIS_ALLOC_SIZE uint16 = 32
	PROJ_ALLOC_SIZE   uint16 = 256

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

	// Player ship
	ship = player{
		pos: pixel.V(WINX/2, 30),
		vel: pixel.ZV,
		acc: 1.1,
		frc: 1 - 0.08,
		hitbox: circularHitbox{
			radius: 12,
			offset: pixel.ZV,
		},
		power:  255,
		sprite: loadSpritesheet("assets/ship-spritesheet.png", pixel.V(13, 18), 3, 7),
	}

	gunCooldown         uint16
	reloadDelay         uint32 = 4
	currentRollCooldown uint16

	// Border values (top, bottom, sides)
	windowBorder = [3]float64{0, 0, 0}
	forceBorder  = [3]float64{400, 150, 30}
	spawnBorder  = [3]float64{-300, -50, -100}

	// Allocation
	projectiles []projectile = make([]projectile, 0, PROJ_ALLOC_SIZE)
	enemies     []enemy      = make([]enemy, 0, ENEMY_ALLOC_SIZE)
	fragments   []fragment   = make([]fragment, 0, ENEMY_ALLOC_SIZE)
	shipTrail   []trailPart  = make([]trailPart, 0, 64)

	// Spritesheets
	projectileSheet pixel.Picture = loadPicture("assets/projectile-spritesheet.png")
	trailSheet      pixel.Picture = loadPicture("assets/trail.png")
	starSheet       pixel.Picture = loadPicture("assets/star-spritesheet.png")

	// Batches
	trailBatch      *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, trailSheet)
	projectileBatch *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, projectileSheet)
	fragmentBatch   *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, fragSheet)

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

func divFloat(n uint16, d uint16) float64 { return float64(n) / float64(d) }

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
	var lenProj uint16 = uint16(len(projectiles))
	for i := lenProj - 1; i < lenProj; i-- {
		if projectiles[i].friendly == friendliness && projectiles[i].pos.Sub(point).Len() < radius {
			inside = append(inside, i)
			count++
		}
	}

	return inside, count
}
