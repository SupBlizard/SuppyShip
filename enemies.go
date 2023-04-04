package main

import "github.com/faiface/pixel"

const ENEMY_ALLOC_SIZE int = 16

var enemies [ENEMY_ALLOC_SIZE]enemy

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
			radius: 16,
			offset: pixel.ZV,
		},
		health: 5,
		sprite: loadSpritesheet("assets/asteroid-spritesheet.png", pixel.V(16, 16), 2),
		name:   "Asteroid",
		id:     0,
	},
}

func updateEnemies() {
	for i := 0; i < ENEMY_ALLOC_SIZE; i++ {
		if !enemies[i].loaded {
			continue
		}

		switch enemies[i].id {
		case 0:
			asteroid(&enemies[i])
		}
	}
}

// AI Functions
func asteroid(ast *enemy) {
	bullets, count := projectilesWithinRadius(ast.phys.pos, ast.hitbox.radius, true)
	if count > 0 {
		ast.health -= count
		// Unload projectiles that collided
		for _, idx := range bullets {
			projectiles[idx].loaded = false
		}
	}

	if ast.health < 1 {
		ast.loaded = false
		return
	}

	ast.sprite.current = ast.health - 1
	drawSprite(&ast.sprite, ast.phys.pos)
}
