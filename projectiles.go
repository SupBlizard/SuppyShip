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
var projectilesLoaded bool

// Projectile Rendering related
var projectileSheet pixel.Picture = loadPicture("assets/projectile-spritesheet.png")
var projectileBatch *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, projectileSheet)
var projSprSize pixel.Vec = pixel.V(6, 16)

// Structs
type projectile struct {
	id         int
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
	vel: pixel.V(0, 12),
	acc: 0,
	frc: 0,
}

var projectileTypes = [4]projectile{
	{
		id:         0,
		name:       "Ship Bullet",
		phys:       shipBulletPhys,
		loaded:     true,
		friendly:   true,
		isAltered:  0,
		cycleSpeed: 15,
		scale:      1,
	},
	{
		id:         1,
		name:       "Onyx Bullet",
		phys:       shipBulletPhys,
		loaded:     true,
		friendly:   true,
		isAltered:  0,
		cycleSpeed: 15,
		scale:      3,
	},
	{
		id:         2,
		name:       "Debris",
		phys:       shipBulletPhys,
		loaded:     true,
		friendly:   false,
		isAltered:  0,
		cycleSpeed: 15,
		scale:      4,
	},
}

// Store the projectile sprite positions on the respective projectiles
func loadProjectileSpritePos() {
	for y := projectileSheet.Bounds().Min.Y; y < projectileSheet.Bounds().Max.Y; y += projSprSize.Y {
		for x := projectileSheet.Bounds().Min.X; x < projectileSheet.Bounds().Max.X; x += projSprSize.X {
			projectileTypes[int(y/projSprSize.Y)].spritesPos[int(x/projSprSize.X)] = pixel.R(x, y, x+projSprSize.X, y+projSprSize.Y)
		}
	}
}

// Spawn a new projectile if one can be loaded
func spawnProjectile(projType int, pos pixel.Vec, vel pixel.Vec) {
	// Find an unloaded projectile slot
	for i := 0; i < BULLET_ALLOC_SIZE; i++ {
		// Ignore loaded slots
		if projectiles[i].loaded {
			continue
		}

		projectiles[i] = projectileTypes[projType]
		projectiles[i].phys.pos = pos
		projectiles[i].phys.vel = vel

		// Return (otherwise all projectiles will be spawned)
		return
	}
}

// Fire a new bullet
func fireBullet(shipPos pixel.Vec) {
	// Check if an Onyx bullet should be created
	indicies, count := bulletsWithinRadius(shipPos, ONYX_CLUSTER_RADIUS)
	if count >= ONYX_CLUSTER_REQUIREMENT {
		// Unload projectiles used
		for _, idx := range indicies {
			projectiles[idx].loaded = false
		}

		// Spawn Onyx bullet
		spawnProjectile(1, shipPos.Add(shipBulletPhys.pos), shipBulletPhys.vel)
		gunCooldown = ONYX_COOLDOWN
	} else {
		spawnProjectile(0, shipPos.Add(shipBulletPhys.pos), shipBulletPhys.vel)
	}

	projectilesLoaded = true
}

// Return all of the bullets within a certain radius around a point
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

// Update the state of each bullet for one frame
func updateProjectiles() {
	if gunCooldown > 0 {
		gunCooldown--
	}
	if !projectilesLoaded {
		return
	}

	projectilesLoaded = false
	for i := 0; i < BULLET_ALLOC_SIZE; i++ {
		if !projectiles[i].loaded {
			continue
		} else if inBounds(projectiles[i].phys.pos, zeroBorder) != pixel.ZV {
			projectiles[i].loaded = false
			continue
		}

		// Count as loaded projectile
		projectilesLoaded = true

		// Animation cycle speed
		if skipFrames(projectiles[i].cycleSpeed) {
			if projectiles[i].isAltered == 0 {
				projectiles[i].isAltered = 1
			} else {
				projectiles[i].isAltered = 0
			}
		}

		// Add velocity to pos
		projectiles[i].phys.pos = projectiles[i].phys.pos.Add(projectiles[i].phys.vel)

		// Draw projectile
		projectile := pixel.NewSprite(projectileSheet, projectiles[i].spritesPos[projectiles[0].isAltered])
		projectile.Draw(projectileBatch, pixel.IM.Scaled(pixel.ZV, projectiles[i].scale).Moved(projectiles[i].phys.pos))

	}

	// Draw all batched projectiles to the window
	projectileBatch.Draw(win)
	projectileBatch.Clear()
}
