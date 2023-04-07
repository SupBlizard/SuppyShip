package main

import (
	"github.com/faiface/pixel"
)

const ENEMY_ALLOC_SIZE uint8 = 16

var (
	enemies       [ENEMY_ALLOC_SIZE]enemy
	loadedEnemies []uint8 = make([]uint8, 0, ENEMY_ALLOC_SIZE)
)

type enemy struct {
	phys   physObj
	loaded bool
	hitbox circularHitbox
	health uint16
	sprite spriteSheet
	name   string
	id     int
}

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
		sprite: loadSpritesheet("assets/asteroid-spritesheet.png", pixel.V(16, 16), 3),
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
	bullets, count := projectilesWithinRadius(ast.phys.pos, ast.hitbox.radius, true)
	if count > 0 {

		if ast.health <= count {
			// Unload enemy
			unloadEnemy(index)
		} else {
			ast.health -= count
		}

		// Unload projectiles that collided
		for _, idx := range bullets {
			unloadProjectile(idx)
		}
	}
	// TODO: Make sprite stages dynamic to health
	ast.sprite.current = ast.health - 1
	drawSprite(&ast.sprite, ast.phys.pos)
}
