package main

import (
	"math"

	"github.com/faiface/pixel"
)

var enemyTypes = []enemy{
	{
		pos: pixel.ZV,
		vel: pixel.ZV,
		acc: 0,
		frc: 0,
		hitbox: circularHitbox{
			radius: 25,
			offset: pixel.ZV,
		},
		health:    10,
		maxHealth: 10,
		sprite:    loadSpritesheet("assets/asteroid-spritesheet.png", pixel.V(16, 16), 3, 30),
		name:      "Asteroid",
		id:        0,
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
		switch enemies[i].id {
		case 0:
			asteroid(&enemies[i], uint16(i))
		}
	}
}

// AI Functions
func asteroid(ast *enemy, idx uint16) {
	// Despawn enemy if it leaves the spawn border
	if inBounds(ast.pos, spawnBorder) != pixel.ZV {
		unloadEnemy(idx)
		return
	}
	enemyHitbox(ast, idx)
	drawSprite(&ast.sprite, ast.pos, uint16(
		math.Round(divFloat(ast.health, ast.maxHealth)*float64(len(ast.sprite.sheet)/int(ast.sprite.cycleNumber)-1))))
}

// Process enemy hitbox
func enemyHitbox(e *enemy, idx uint16) {
	proj, count := projectilesInRadius(e.pos, e.hitbox.radius, true)
	if count > 0 {
		if e.health <= count {
			unloadEnemy(idx)
		} else {
			e.health -= count
		}

		// Unload projectiles that collided
		unloadMany(proj)
	}
}
