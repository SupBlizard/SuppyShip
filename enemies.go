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
	slot := uint8(0)
	for ; slot < ENEMY_ALLOC_SIZE; slot++ {
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
func unloadEnemy(idx uint8) {
	enemies[idx].loaded = false
	for i := 0; i < len(loadedEnemies); i++ {
		if loadedEnemies[i] == idx {
			loadedEnemies = append(loadedEnemies[:i], loadedEnemies[i+1:]...)
			return
		}
	}
}

func updateEnemies() {
	for _, i := range loadedEnemies {
		switch enemies[i].id {
		case 0:
			asteroid(&enemies[i], i)
		}
	}
}

// AI Functions
func asteroid(ast *enemy, index uint8) {
	// Despawn enemy if it leaves the spawn border
	if inBounds(ast.phys.pos, spawnBorder) != pixel.ZV {
		unloadEnemy(index)
		return
	}

	bullets, count := projectilesInRadius(ast.phys.pos, ast.hitbox.radius, true)
	if count > 0 {

		if ast.health <= count {
			// Unload enemy
			unloadEnemy(index)
		} else {
			ast.health -= count
		}

		// Unload projectiles that collided
		for _, loadID := range bullets {
			unloadProjectile(loadID)
		}
	}
	// TODO: Make sprite stages dynamic to health
	drawSprite(&ast.sprite, ast.phys.pos, ast.health-1)
}
