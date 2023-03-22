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

// Structs
type projectile struct {
	name        string
	phys        physObj
	loaded      bool
	friendly    bool
	isAltered   bool
	cycleSpeed  int
	scale       float64
	spritesheet [2]*pixel.Sprite
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
		isAltered:  false,
		cycleSpeed: 15,
		scale:      1,
	},
	{
		name:       "Onyx Bullet",
		phys:       shipBulletPhys,
		loaded:     true,
		friendly:   true,
		isAltered:  false,
		cycleSpeed: 15,
		scale:      3,
	},
	{
		name:       "Debris",
		phys:       shipBulletPhys,
		loaded:     true,
		friendly:   false,
		isAltered:  false,
		cycleSpeed: 15,
		scale:      4,
	},
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

			// Draw projectile
			drawProjectile(projectiles[i])
		}
	}
}

func drawProjectile(proj projectile) {
	if frameCount%proj.cycleSpeed == 0 {
		proj.isAltered = !proj.isAltered
	}

	spriteIndex := 0
	if proj.isAltered {
		spriteIndex = 1
	}
	proj.spritesheet[spriteIndex].Draw(win, pixel.IM.Scaled(pixel.ZV, proj.scale).Moved(proj.phys.pos))
}

func loadProjectileSprites() {
	size := pixel.V(6, 16)
	img, err := loadPicture("assets/projectile-spritesheet.png")
	if err != nil {
		panic(err)
	}

	// set to min after
	for y := 0.0; y < img.Bounds().Max.Y; y += size.Y {
		for x := 0.0; x < img.Bounds().Max.X; x += size.X {
			projectileTypes[int(y/size.Y)].spritesheet[int(x/size.X)] = pixel.NewSprite(img, pixel.R(x, y, x+size.X, y+size.Y))
		}
	}
}
