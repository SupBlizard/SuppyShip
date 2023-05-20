package main

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
)

var ship = player{
	pos:  pixel.V(WINX/2, 40),
	vel:  pixel.ZV,
	dir:  0,
	roll: 0,
	acc:  1.1,
	frc:  1 - 0.08,

	power:    30,
	recharge: false,
	alive:    true,
	heat:     0,
	reload:   5,

	hitbox: circularHitbox{radius: 12, offset: pixel.ZV},
	sprite: loadSpritesheet("assets/ship-spritesheet.png", pixel.V(13, 18), 3, 7),
	frag:   fragInfo{ID: 0, frags: 3, power: 0.5, radius: 5, scale: 3},
	shield: shipShield{
		prot:       0,
		protLength: 60,
		sprite:     loadSpritesheet("assets/shield.png", pixel.V(34, 34), 2, 10),
	},
}

func updateShip() {
	// Give velocity a minimum limit or apply friction to the velocity if there is any
	if ship.vel.Len() <= 0.01 {
		ship.vel = pixel.ZV
	} else {
		ship.vel = ship.vel.Scaled(ship.frc)
	}

	// Handle rolling after friction being applied
	handleRolling()

	// Update shield state
	updateShield()

	// Add new velocity if there is input
	if input.dir != pixel.ZV {
		ship.vel = ship.vel.Add(input.dir)
	}

	// Enforce soft boundary on ship
	if borderCollisions := inBounds(ship.pos, forceBorder); borderCollisions != pixel.ZV {
		var borderDepth float64
		var globalAccIdx int

		if borderCollisions.Y == -1 {
			borderDepth = findBorderDepth(WINY-ship.pos.Y, forceBorder[0])
			globalAccIdx = 0
		} else if borderCollisions.Y == 1 {
			borderDepth = findBorderDepth(ship.pos.Y, forceBorder[1])
			globalAccIdx = 1
		}

		counterAcceleration := ship.acc * BOUNDARY_STRENGTH
		globalVelocity -= (borderDepth * borderCollisions.Y * math.Pow(globalAcc[globalAccIdx], 2))
		ship.vel.Y += counterAcceleration * borderDepth * borderCollisions.Y

		if borderCollisions.X == -1 {
			ship.vel.X -= counterAcceleration * findBorderDepth(WINX-ship.pos.X, forceBorder[2])
		} else if borderCollisions.X == 1 {
			ship.vel.X += counterAcceleration * findBorderDepth(ship.pos.X, forceBorder[2])
		}
	}

	// Add new velocity to the position
	if ship.vel.Len() != 0 {
		ship.pos = ship.pos.Add(ship.vel)
	}
}

func collidingWithShip(obj *enemy) bool {
	diffVec := ship.pos.Sub(obj.pos)
	if diffVec.Len() < obj.hitbox.radius+ship.hitbox.radius {
		// Bounce ship off
		ship.vel = diffVec.Scaled(ship.vel.Len()/diffVec.Len() + 0.05)
		return true
	}
	return false
}

func hitShip() {
	if ship.shield.prot > 0 {
		return
	}
	if ship.shield.active {
		ship.power = 0
		ship.shield.prot = ship.shield.protLength
		return
	}
	fragmentObject(&ship.frag, []uint8{0, 1, 2}, ship.pos, ship.vel)
	ship.alive = false
	ship.pos = pixel.V(WINX/2, -200)
}

func updateShield() {
	if ship.power == 0xFF {
		ship.hitbox.radius = SHIELD_RADIUS
		ship.shield.active = true
	} else {
		ship.hitbox.radius = SHIP_RADIUS
		ship.shield.active = false
	}
}

// Draw ship to the screen
func drawShip() {
	var spriteID uint16 = 6
	if ship.roll == 0 {
		if input.dir.Y != 0 {
			if input.dir.Y < 0 {
				spriteID = 5
			} else if globalVelocity < 0.5 {
				spriteID = 7
				shipTrail = append(shipTrail, trailPart{pos: ship.pos, mask: color.RGBA{255, 255, 255, 255}})
			} else {
				spriteID = 8
				shipTrail = append(shipTrail, trailPart{pos: ship.pos, mask: color.RGBA{255, 255, 255, 255}})
			}
		}

		if math.Abs(input.dir.X) > AXIS_DEADZONE {
			spriteID = 4
			if input.dir.X > 0 {
				spriteID = 0
			}
		}
	} else {
		// (ship.dir*-1*4)) was the if statement for offset (offset 0 and 4)
		spriteID = uint16(int16(ship.roll/(ROLL_COOLDOWN/4))*ship.dir*-1+(ship.dir*4)) % ROLL_SPRITE_NUMBER
	}

	drawSprite(&ship.sprite, ship.pos, 0, spriteID)
}

// Fire bullet
func fireBullet(shipPos pixel.Vec) {
	bullets := projectilesInRadius(shipPos, ONYX_CLUSTER_RADIUS, true)

	// Check if an Onyx bullet should be created
	if uint16(len(bullets)) >= ONYX_CLUSTER_REQUIREMENT {
		unloadMany(bullets)

		// Spawn Onyx bullet
		loadProjectile(1, shipPos.Add(projectileTypes[1].pos), projectileTypes[1].vel)
		ship.heat = ONYX_COOLDOWN
		return
	}

	loadProjectile(0, shipPos.Add(projectileTypes[0].pos), projectileTypes[0].vel)
}

func unloadTrailPart(ID int) {
	shipTrail[ID] = shipTrail[len(shipTrail)-1]
	shipTrail = shipTrail[:len(shipTrail)-1]
}

func updateShipTrail(shipPos pixel.Vec) {
	if len(shipTrail) == 0 {
		return
	}

	for i := 0; i < len(shipTrail); i++ {
		shipTrail[i].mask.A -= 20
		if shipTrail[i].mask.A < 20 || math.Abs(shipTrail[i].pos.X-shipPos.X) > 6 {
			unloadTrailPart(i)
			continue
		}

		shipTrail[i].pos = shipTrail[i].pos.Sub(pixel.V(0, SHIPTRAIL_LENGTH+(globalVelocity*SHIPTRAIL_ACC)-ship.vel.Y))
		scale := 2.5 * (float64(shipTrail[i].mask.A) / 0xFF)

		// Draw trail
		pixel.NewSprite(trailSheet, trailSheet.Bounds()).Draw(
			trailBatch, pixel.IM.Scaled(pixel.ZV, scale).Moved(shipTrail[i].pos))
	}

	trailBatch.Draw(win)
	trailBatch.Clear()
}

func handleRolling() {
	dir := signbit(ship.vel.X)
	if ship.roll == 0 {
		if input.roll && math.Abs(ship.vel.X) > 0.5 {
			ship.roll = ROLL_COOLDOWN
			ship.dir = int16(dir)
			ship.vel.X += 10 * dir
		}
	} else {
		// Dampen control on the X axis
		input.dir.X /= 2
		ship.roll--
	}
}
