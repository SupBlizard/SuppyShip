package main

import (
	"github.com/faiface/pixel"
)

var enemyTypes = []enemy{
	{
		pos: pixel.ZV,
		vel: pixel.ZV,
		acc: 0,
		frc: 0,
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

// Load a new projectile if there is space
func loadEnemy(enemyType uint8, pos pixel.Vec, vel pixel.Vec) {
	if enemyLen := uint16(len(enemies)); enemyLen < ENEMY_ALLOC_SIZE {
		enemies = append(enemies, enemyTypes[enemyType])
		enemies[enemyLen].pos = pos
		enemies[enemyLen].vel = vel
	}
}

// Unload projectiles
func unloadEnemy(idx uint16) {
	enemies[idx] = enemies[len(enemies)-1]
	enemies = enemies[:len(enemies)-1]
}

func updateEnemies() {
	for i := len(enemies) - 1; i > -1; i-- {
		println(enemies)
		switch enemies[i].id {
		case 0:
			asteroid(&enemies[i], uint16(i))
		}
	}
}

// AI Functions
func asteroid(ast *enemy, index uint16) {
	// Despawn enemy if it leaves the spawn border
	if inBounds(ast.pos, spawnBorder) != pixel.ZV {
		unloadEnemy(index)
		return
	}

	bullets, count := projectilesInRadius(ast.pos, ast.hitbox.radius, true)
	if count > 0 {
		if ast.health <= count {
			unloadEnemy(index)
		} else {
			ast.health -= count
		}

		// Unload projectiles used
		unloadMany(bullets)
	}
	// TODO: Make sprite stages dynamic to health
	drawSprite(&ast.sprite, ast.pos, ast.health-1)
}
