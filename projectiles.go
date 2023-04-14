package main

import (
	"github.com/faiface/pixel"
)

var projectileTypes = [4]projectile{
	{
		id:       0,
		name:     "Ship Bullet",
		pos:      pixel.V(0, 10),
		vel:      pixel.V(0, 12),
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
		pos:      pixel.V(0, 10),
		vel:      pixel.V(0, 12),
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
		pos:      pixel.V(0, 10),
		vel:      pixel.V(0, 12),
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
	const SPR_SIZE_X float64 = 6
	const SPR_SIZE_Y float64 = 16
	for y := projectileSheet.Bounds().Min.Y; y < projectileSheet.Bounds().Max.Y; y += SPR_SIZE_Y {
		for x := projectileSheet.Bounds().Min.X; x < projectileSheet.Bounds().Max.X; x += SPR_SIZE_X {
			projectileTypes[int(y/SPR_SIZE_Y)].sprite.pos[int(x/SPR_SIZE_X)] = pixel.R(x, y, x+SPR_SIZE_X, y+SPR_SIZE_Y)
		}
	}
}

// Load a new projectile if there is space
func loadProjectile(projType uint8, pos pixel.Vec, vel pixel.Vec) {
	// Find first free slot
	for slot := uint16(0); slot < BULLET_ALLOC_SIZE; slot++ {
		if projectiles[slot].loaded {
			continue
		}
		// Fill in slot
		projectiles[slot] = projectileTypes[projType]
		projectiles[slot].pos = pos
		projectiles[slot].vel = vel

		// Add projectile to the loaded list
		loadedProjectiles = append(loadedProjectiles, slot)

		return
	}
}

// Unload projectiles
func unloadProjectile(idx uint16) {
	for i := 0; i < len(loadedProjectiles); i++ {
		if loadedProjectiles[i] == idx {
			projectiles[loadedProjectiles[i]].loaded = false
			loadedProjectiles[i] = loadedProjectiles[len(loadedProjectiles)-1]
			loadedProjectiles = loadedProjectiles[:len(loadedProjectiles)-1]
			return
		}
	}
}

// Fire a new bullet
func fireBullet(shipPos pixel.Vec) {
	// Check if an Onyx bullet should be created
	bullets, count := projectilesInRadius(shipPos, ONYX_CLUSTER_RADIUS, true)
	if count >= ONYX_CLUSTER_REQUIREMENT {
		// Unload projectiles used
		for _, projID := range bullets {
			unloadProjectile(projID)
		}

		// Spawn Onyx bullet
		loadProjectile(1, shipPos.Add(projectileTypes[1].pos), projectileTypes[1].vel)
		gunCooldown = ONYX_COOLDOWN
	} else {
		loadProjectile(0, shipPos.Add(projectileTypes[0].pos), projectileTypes[0].vel)
	}
}

// Update the state of each bullet for one frame
func updateProjectiles() {
	if gunCooldown > 0 {
		gunCooldown--
	}

	// Loop through loaded indexes
	for _, i := range loadedProjectiles {
		// Unload out of bounds projectiles
		if inBounds(projectiles[i].pos, windowBorder) != pixel.ZV {
			unloadProjectile(i)
			i -= 1
			continue
		}

		// Animation cycle speed
		if skipFrames(projectiles[i].sprite.cycleSpeed) {
			if projectiles[i].sprite.cycle == 0 {
				projectiles[i].sprite.cycle = 1
			} else {
				projectiles[i].sprite.cycle = 0
			}
		}

		// Add velocity to pos
		projectiles[i].pos = projectiles[i].pos.Add(projectiles[i].vel)

		// Draw projectile
		pixel.NewSprite(projectileSheet, projectiles[i].sprite.pos[projectiles[0].sprite.cycle]).Draw(
			projectileBatch, pixel.IM.Scaled(pixel.ZV, projectiles[i].sprite.scale).Moved(projectiles[i].pos),
		)

	}

	// Draw all batched projectiles to the window
	projectileBatch.Draw(win)
	projectileBatch.Clear()
}
