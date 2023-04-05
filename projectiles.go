package main

import (
	"github.com/faiface/pixel"
)

const BULLET_ALLOC_SIZE uint16 = 256
const ONYX_CLUSTER_REQUIREMENT uint16 = 7
const ONYX_CLUSTER_RADIUS float64 = 30
const ONYX_COOLDOWN int = 60

var (
	reloadDelay int = 4
	gunCooldown     = 0

	// Projectile Allocation
	projectiles       [BULLET_ALLOC_SIZE]projectile
	loadedProjectiles []uint16 = make([]uint16, 0, BULLET_ALLOC_SIZE)

	// Projectile Rendering related
	projectileSheet pixel.Picture = loadPicture("assets/projectile-spritesheet.png")
	projectileBatch *pixel.Batch  = pixel.NewBatch(&pixel.TrianglesData{}, projectileSheet)
	projSprSize     pixel.Vec     = pixel.V(6, 16)
)

// Structs
type projectile struct {
	id       int
	name     string
	phys     physObj
	loaded   bool
	friendly bool
	sprite   projectileSprite
}

type projectileSprite struct {
	cycle      uint8
	cycleSpeed int
	scale      float64
	pos        [2]pixel.Rect
}

var shipBulletPhys physObj = physObj{
	pos: pixel.V(0, 10),
	vel: pixel.V(0, 0.5),
	acc: 0,
	frc: 0,
}

var projectileTypes = [4]projectile{
	{
		id:       0,
		name:     "Ship Bullet",
		phys:     shipBulletPhys,
		loaded:   true,
		friendly: true,
		sprite: projectileSprite{
			cycle:      0,
			cycleSpeed: 15,
			scale:      1,
		},
	},
	{
		id:       1,
		name:     "Onyx Bullet",
		phys:     shipBulletPhys,
		loaded:   true,
		friendly: true,
		sprite: projectileSprite{
			cycle:      0,
			cycleSpeed: 15,
			scale:      3,
		},
	},
	{
		id:       2,
		name:     "Debris",
		phys:     shipBulletPhys,
		loaded:   true,
		friendly: false,
		sprite: projectileSprite{
			cycle:      0,
			cycleSpeed: 15,
			scale:      4,
		},
	},
}

// Store the projectile sprite positions on the respective projectiles
func loadProjectileSpritePos() {
	for y := projectileSheet.Bounds().Min.Y; y < projectileSheet.Bounds().Max.Y; y += projSprSize.Y {
		for x := projectileSheet.Bounds().Min.X; x < projectileSheet.Bounds().Max.X; x += projSprSize.X {
			projectileTypes[int(y/projSprSize.Y)].sprite.pos[int(x/projSprSize.X)] = pixel.R(x, y, x+projSprSize.X, y+projSprSize.Y)
		}
	}
}

// Load a new projectile if there is space
func loadProjectile(projType int, pos pixel.Vec, vel pixel.Vec) {
	// Find an unloaded projectile slot
	for i := uint16(0); i < BULLET_ALLOC_SIZE; i++ {
		// Ignore loaded slots
		if projectiles[i].loaded {
			continue
		}

		projectiles[i] = projectileTypes[projType]
		projectiles[i].phys.pos = pos
		projectiles[i].phys.vel = vel

		// Add projectile to the loaded list
		loadedProjectiles = append(loadedProjectiles, i)

		// Return (otherwise all projectiles will be spawned)
		return
	}
}

// Unload projectiles
func unloadProjectile(idx uint16) {
	projectiles[idx].loaded = false
	for i := 0; i < len(loadedProjectiles); i++ {
		if loadedProjectiles[i] == idx {
			loadedProjectiles = append(loadedProjectiles[:i], loadedProjectiles[i+1:]...)
			return
		}
	}
}

// Fire a new bullet
func fireBullet(shipPos pixel.Vec) {
	// Check if an Onyx bullet should be created
	indicies, count := projectilesWithinRadius(shipPos, ONYX_CLUSTER_RADIUS, true)
	if count >= ONYX_CLUSTER_REQUIREMENT {
		// Unload projectiles used
		for _, idx := range indicies {
			unloadProjectile(idx)
		}

		// Spawn Onyx bullet
		loadProjectile(1, shipPos.Add(shipBulletPhys.pos), shipBulletPhys.vel)
		gunCooldown = ONYX_COOLDOWN
	} else {
		loadProjectile(0, shipPos.Add(shipBulletPhys.pos), shipBulletPhys.vel)
	}
}

// Return all of the projectiles within a certain radius around a point
func projectilesWithinRadius(point pixel.Vec, radius float64, friendliness bool) ([]uint16, uint16) {
	var inside []uint16
	var count uint16 = 0

	var idx uint16
	for i := uint16(0); i < uint16(len(loadedProjectiles)); i++ {
		idx = loadedProjectiles[i]

		// Check if the projectile is of a specific friendliness and within a radius
		if projectiles[idx].friendly == friendliness && projectiles[idx].phys.pos.Sub(point).Len() < radius {
			inside = append(inside, i)
			count++
		}
	}
	return inside, count
}

// Update the state of each bullet for one frame
func updateProjectiles() {
	if gunCooldown > 0 {
		gunCooldown--
	}

	var idx uint16
	for i := 0; i < len(loadedProjectiles); i++ {
		idx = loadedProjectiles[i]

		// Unload out of bounds projectiles
		if inBounds(projectiles[idx].phys.pos, zeroBorder) != pixel.ZV {
			unloadProjectile(idx)
			continue
		}

		// Animation cycle speed
		if skipFrames(projectiles[idx].sprite.cycleSpeed) {
			if projectiles[idx].sprite.cycle == 0 {
				projectiles[idx].sprite.cycle = 1
			} else {
				projectiles[idx].sprite.cycle = 0
			}
		}

		// Add velocity to pos
		projectiles[idx].phys.pos = projectiles[idx].phys.pos.Add(projectiles[idx].phys.vel)

		// Draw projectile
		projectile := pixel.NewSprite(projectileSheet, projectiles[idx].sprite.pos[projectiles[0].sprite.cycle])
		projectile.Draw(projectileBatch, pixel.IM.Scaled(pixel.ZV, projectiles[idx].sprite.scale).Moved(projectiles[idx].phys.pos))

	}

	// Draw all batched projectiles to the window
	projectileBatch.Draw(win)
	projectileBatch.Clear()
}
