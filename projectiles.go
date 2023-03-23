package main

import (
	"github.com/faiface/pixel"
)

const BULLET_ALLOC_SIZE int = 256
const ONYX_CLUSTER_REQUIREMENT int = 7
const ONYX_CLUSTER_RADIUS float64 = 30
const ONYX_COOLDOWN int = 60

var reloadDelay int = 4
var gunCooldown = 0

// Projectile Allocation Array
var projectiles [BULLET_ALLOC_SIZE]projectile

// Projectile Rendering related
var projectileSheet pixel.Picture = loadPicture("assets/projectile-spritesheet.png")
var projectileBatch *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, projectileSheet)
var projSprSize pixel.Vec = pixel.V(6, 16)

// Structs
type projectile struct {
	name       string
	phys       physObj
	loaded     bool
	friendly   bool
	isAltered  uint8
	cycleSpeed int
	scale      float64
	spritesPos [2]pixel.Rect
}

var shipBulletPhys physObj = physObj{
	pos: pixel.V(0, 10),
	vel: pixel.V(0, 10),
	acc: 0,
	frc: 0,
}

var projectileTypes = [4]projectile{
	{
		name:       "Bullet",
		phys:       shipBulletPhys,
		loaded:     true,
		friendly:   true,
		isAltered:  0,
		cycleSpeed: 15,
		scale:      1,
	},
	{
		name:       "Onyx Bullet",
		phys:       shipBulletPhys,
		loaded:     true,
		friendly:   true,
		isAltered:  0,
		cycleSpeed: 15,
		scale:      3,
	},
	{
		name:       "Debris",
		phys:       shipBulletPhys,
		loaded:     true,
		friendly:   false,
		isAltered:  0,
		cycleSpeed: 15,
		scale:      4,
	},
}

func loadProjectileSpritePos() {
	for y := projectileSheet.Bounds().Min.Y; y < projectileSheet.Bounds().Max.Y; y += projSprSize.Y {
		for x := projectileSheet.Bounds().Min.X; x < projectileSheet.Bounds().Max.X; x += projSprSize.X {
			projectileTypes[int(y/projSprSize.Y)].spritesPos[int(x/projSprSize.X)] = pixel.R(x, y, x+projSprSize.X, y+projSprSize.Y)
		}
	}
}

// [Create a bullet if a slot is free]
func createBullet(shipPos pixel.Vec) {
	for i := 0; i < BULLET_ALLOC_SIZE; i++ {
		if !projectiles[i].loaded {
			projectiles[i] = projectileTypes[0]
			projectiles[i].phys.pos = projectiles[i].phys.pos.Add(shipPos)

			indicies, count := bulletsWithinRadius(projectiles[i].phys.pos, ONYX_CLUSTER_RADIUS)
			if count >= ONYX_CLUSTER_REQUIREMENT {
				unloadProjectiles(indicies)
				bulletPos := projectiles[i].phys.pos
				projectiles[i] = projectileTypes[1]
				projectiles[i].phys.pos = bulletPos

				gunCooldown = ONYX_COOLDOWN
			}
			return
		}
	}
}

// [Unload all selected projectiles]
func unloadProjectiles(selected []int) {
	for i := 0; i < len(selected); i++ {
		projectiles[selected[i]].loaded = false
	}
}

// [Return all of the bullets within a certain radius around a point]
func bulletsWithinRadius(point pixel.Vec, radius float64) ([]int, int) {
	var insideRadius []int
	var projectileCount int = 0
	for i := 0; i < BULLET_ALLOC_SIZE; i++ {
		if projectiles[i].loaded && projectiles[i].phys.pos.Sub(point).Len() < radius {
			insideRadius = append(insideRadius, i)
			projectileCount++
		}
	}
	return insideRadius, projectileCount
}

// [Update the state of each bullet for one frame]
func updateBullets() {
	// Update bullets
	for i := 0; i < BULLET_ALLOC_SIZE; i++ {
		if !projectiles[i].loaded {
			continue
		}
		if inBounds(projectiles[i].phys.pos, WINSIZE, NULL_BOUNDARY_RANGE) != pixel.ZV {
			projectiles[i].loaded = false
		} else {
			projectiles[i].phys.pos = projectiles[i].phys.pos.Add(projectiles[i].phys.vel)

		}
	}

	// Draw projectiles
	drawProjectiles()
}

func drawProjectiles() {

	for i := 0; i < BULLET_ALLOC_SIZE; i++ {
		if !projectiles[i].loaded {
			continue
		}

		if frameCount%projectiles[i].cycleSpeed == 0 {
			if projectiles[i].isAltered == 0 {
				projectiles[i].isAltered = 1
			} else {
				projectiles[i].isAltered = 0
			}
		}

		projectile := pixel.NewSprite(projectileSheet, projectiles[i].spritesPos[projectiles[0].isAltered])
		projectile.Draw(projectileBatch, pixel.IM.Scaled(pixel.ZV, projectiles[i].scale).Moved(projectiles[i].phys.pos))
	}

	// Draw all projectiles
	projectileBatch.Draw(win)
	projectileBatch.Clear()
}
