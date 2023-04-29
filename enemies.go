package main

import (
	"math"

	"github.com/faiface/pixel"
)

var enemyTypes = []enemy{
	{
		id:        0,
		rotVel:    0.01,
		health:    10,
		maxHealth: 10,
		hitbox:    circularHitbox{radius: 25, offset: pixel.ZV},
		sprite:    loadSpritesheet("assets/asteroid-spritesheet.png", pixel.V(16, 16), 3, 30),
		frag: fragInfo{
			ID:     1,
			frags:  8,
			power:  2,
			radius: 15,
			scale:  3,
		},
	},
	{
		id:        1,
		health:    10,
		maxHealth: 10,
		hitbox:    circularHitbox{radius: 25, offset: pixel.ZV},
		sprite:    loadSpritesheet("assets/eye-spritesheet.png", pixel.V(16, 16), 3, 30),
		frag: fragInfo{
			ID:     2,
			frags:  8,
			power:  2,
			radius: 15,
			scale:  3,
		},
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
	fragmentObject(&enemies[idx].frag, []uint8{}, enemies[idx].pos, enemies[idx].vel)
	enemies[idx] = enemies[len(enemies)-1]
	enemies = enemies[:len(enemies)-1]
}

// Process enemy hitbox
func enemyHitbox(enemyID uint16) uint16 {
	bullets := projectilesInRadius(enemies[enemyID].pos, enemies[enemyID].hitbox.radius, true)
	if len(bullets) == 0 {
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

		// Check if colliding with ship
		if ship.alive && collidingWithShip(&enemies[i]) {
			hitShip()
		}

		// Calculate enemy health
		var damage uint16 = enemyHitbox(i)
		if enemies[i].health <= damage {
			unloadEnemy(i)
			continue
		} else if damage != 0 {
			enemies[i].health -= damage
		}

		// Add enermy rotation
		enemies[i].rot += enemies[i].rotVel

		// Process custom enemy code
		switch enemies[i].id {
		case 0:
			asteroid(i)
		case 1:
			eye(i)
		}
	}
}

// AI Functions
func asteroid(enemyID uint16) {
	var ast *enemy = &enemies[enemyID]

	ast.pos = ast.pos.Add(ast.vel)

	drawSprite(&ast.sprite, ast.pos, ast.rot, uint16(
		math.Round(divFloat(ast.health, ast.maxHealth)*float64(len(ast.sprite.sheet)/int(ast.sprite.cycleNumber)-1))))
}

func eye(enemyID uint16) {
	var eye *enemy = &enemies[enemyID]

	eye.pos = eye.pos.Add(eye.vel)

	sight := ship.pos.Sub(eye.pos)
	shouldLook := math.Acos(sight.X/sight.Len()) * signbit(sight.Y)
	eye.rotVel = (shouldLook - eye.rot)
	if (shouldLook - eye.rot) < AXIS_DEADZONE {
		eye.rotVel = (shouldLook - eye.rot)
	}
	drawSprite(&eye.sprite, eye.pos, eye.rot, uint16(
		math.Round(divFloat(eye.health, eye.maxHealth)*float64(len(eye.sprite.sheet)/int(eye.sprite.cycleNumber)-1))))
}
