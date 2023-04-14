package main

import (
	"github.com/faiface/pixel"
)

var enemyTypes = []enemy{
	{
		phys: physObj{
			pos: pixel.ZV,
			vel: pixel.ZV,
			acc: 0,
			frc: 0,
		},
		loaded: true,
		hitbox: circularHitbox{
			radius: 20,
			offset: pixel.ZV,
		},
		health: 5,
		sprite: loadSpritesheet("assets/asteroid-spritesheet.png", pixel.V(16, 16), 3, 30),
		name:   "Asteroid",
		id:     0,
	},
}

// Load a new enemy if there is space
func loadEnemy(enemyType int, pos pixel.Vec, vel pixel.Vec) {
	// Find first free slot
	for slot := uint16(0); slot < ENEMY_ALLOC_SIZE; slot++ {
		if enemies[slot].loaded {
			continue
		}

		// Fill in slot
		enemies[slot] = enemyTypes[enemyType]
		enemies[slot].phys.pos = pos
		enemies[slot].phys.vel = vel

		// Add enemy to the loaded list
		loadedEnemies = append(loadedEnemies, slot)
		return
	}
}

// Unload enemies
func unloadEnemy(idx uint16) {
	enemies[loadedEnemies[idx]].loaded = false
	loadedEnemies[idx] = loadedEnemies[len(loadedEnemies)-1]
	loadedEnemies = loadedEnemies[:len(loadedEnemies)-1]
}

func updateEnemies() {
	for _, loadID := range loadedEnemies {
		switch enemies[loadedEnemies[loadID]].id {
		case 0:
			asteroid(&enemies[loadedEnemies[loadID]], loadID)
		}
	}
}

// AI Functions
func asteroid(ast *enemy, index uint16) {
	// Despawn enemy if it leaves the spawn border
	if inBounds(ast.phys.pos, spawnBorder) != pixel.ZV {
		unloadEnemy(index)
		return
	}

	bullets, count := projectilesInRadius(ast.phys.pos, ast.hitbox.radius, true)
	if count > 0 {
		if ast.health <= count {
			unloadEnemy(index)
		} else {
			ast.health -= count
		}

		// Find projectile IDs to unload
		var projIDS []uint16 = make([]uint16, 0, 16)
		for _, loadID := range bullets {
			projIDS = append(projIDS, loadedProjectiles[loadID])
		}

		// Unload projectiles used
		for _, projID := range projIDS {
			unloadProjectile(findLoadID(loadedProjectiles, projID))
		}
	}
	// TODO: Make sprite stages dynamic to health
	drawSprite(&ast.sprite, ast.phys.pos, ast.health-1)
}
