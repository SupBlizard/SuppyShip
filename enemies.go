package main

import (
	"math"
	"math/rand"

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
	fragmentObject(1, uint8(rand.Int31()%3), enemies[idx].pos, enemies[idx].vel, 2, 20, 8)
	enemies[idx] = enemies[len(enemies)-1]
	enemies = enemies[:len(enemies)-1]
}

// Process enemy hitbox
func enemyHitbox(enemyID uint16) uint16 {
	bullets, count := projectilesInRadius(enemies[enemyID].pos, enemies[enemyID].hitbox.radius, true)
	if count == 0 {
		return 0
	}

	// Count up total damage
	var damage uint16
	for _, projID := range bullets {
		damage += projectiles[projID].dmg
	}

	// Unload projectiles that collided
	unloadMany(bullets)

	return damage
}

// Update every enemy
func updateEnemies() {
	var lenEnemies uint16 = uint16(len(enemies))
	for i := lenEnemies - 1; i < lenEnemies; i-- {
		// Despawn enemy if it leaves the spawn border
		if inBounds(enemies[i].pos, spawnBorder) != pixel.ZV {
			unloadEnemy(i)
			continue
		}

		// Calculate enemy health
		var damage uint16 = enemyHitbox(i)
		if enemies[i].health <= damage {
			unloadEnemy(i)
			continue
		} else if damage != 0 {
			enemies[i].health -= damage
		}

		// Process custom enemy code
		switch enemies[i].id {
		case 0:
			asteroid(i)
		}
	}
}

// AI Functions
func asteroid(enemyID uint16) {
	var ast *enemy = &enemies[enemyID]

	ast.pos = ast.pos.Add(ast.vel)

	drawSprite(&ast.sprite, ast.pos, uint16(
		math.Round(divFloat(ast.health, ast.maxHealth)*float64(len(ast.sprite.sheet)/int(ast.sprite.cycleNumber)-1))))
}
