package main

import (
	"github.com/faiface/pixel"
)

var projSprSize = pixel.V(6, 16)

var projectileTypes = [2]projectile{
	{
		id:       0,
		pos:      pixel.V(0, 10),
		vel:      pixel.V(0, 12),
		dmg:      1,
		friendly: true,
		sprite: projectileSprite{
			cycle:      0,
			cycleSpeed: 240,
			scale:      1,

			pos: [PROJ_SPRITES]pixel.Rect(batchSpritePos(0, projectileSheet, projSprSize)),
		},
	},
	{
		id:       1,
		pos:      pixel.V(0, 10),
		vel:      pixel.V(0, 12),
		dmg:      5,
		friendly: true,
		sprite: projectileSprite{
			cycle:      0,
			cycleSpeed: 240,
			scale:      3,

			pos: [PROJ_SPRITES]pixel.Rect(batchSpritePos(1, projectileSheet, projSprSize)),
		},
	},
}

// Load a new projectile if there is space
func loadProjectile(projType uint8, pos pixel.Vec, vel pixel.Vec) {
	if projlen := uint16(len(projectiles)); projlen < PROJ_ALLOC_SIZE {
		projectiles = append(projectiles, projectileTypes[projType])
		projectiles[projlen].pos = pos
		projectiles[projlen].vel = vel
	}
}

// Unload projectiles
func unloadProjectile(idx uint16) {
	projectiles[idx] = projectiles[len(projectiles)-1]
	projectiles = projectiles[:len(projectiles)-1]
}

// Unload a whole slice of projectiles (indicies must be in descending order)
func unloadMany(indicies []uint16) {
	for _, idx := range indicies {
		projectiles[idx] = projectiles[len(projectiles)-1]
		projectiles = projectiles[:len(projectiles)-1]
	}
}

// Update the state of each bullet for one frame
func updateProjectiles() {
	if ship.heat > 0 {
		ship.heat--
	}
	if len(projectiles) == 0 {
		return
	}

	// Loop through loaded indexes
	var lenProj uint16 = uint16(len(projectiles))
	for i := lenProj - 1; i < lenProj; i-- {
		// Unload out of bounds projectiles
		if inBounds(projectiles[i].pos, windowBorder) != pixel.ZV {
			unloadProjectile(i)
			continue
		}

		// Animation cycle speed
		if timePassed(projectiles[i].sprite.cycleSpeed) {
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
