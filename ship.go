package main

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
)

// Update everything about the ship
func updateShip() {
	// Give velocity a minimum limit or apply friction to the velocity if there is any
	if ship.vel.Len() <= 0.01 {
		ship.vel = pixel.ZV
	} else {
		ship.vel = ship.vel.Scaled(ship.frc)
	}

	// Add new velocity if there is input
	if input.dir != pixel.ZV {
		ship.vel = ship.vel.Add(input.dir)
	}

	// Handle rolling after friction being applied
	handleRolling()

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
		globalVelocity -= (borderDepth * borderCollisions.Y * (DEFAULT_GLOBAL_VELOCITY * globalAcc[globalAccIdx]))
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

// Draw ship to the screen
func drawShip() {
	var spriteID uint16 = 1
	if currentRollCooldown == 0 {
		if input.dir.Y != 0 {
			if input.dir.Y < 0 {
				spriteID = 0
			} else if globalVelocity < DEFAULT_GLOBAL_VELOCITY+5 {
				spriteID = 2
				shipTrail = append(shipTrail, trailPart{pos: ship.pos.Sub(pixel.V(0, 18)), mask: color.RGBA{255, 255, 255, 255}})
			} else {
				spriteID = 3
				shipTrail = append(shipTrail, trailPart{pos: ship.pos.Sub(pixel.V(0, 18)), mask: color.RGBA{255, 255, 255, 255}})
			}

		}

		if math.Abs(input.dir.X) > AXIS_DEADZONE {
			if input.dir.X > 0 {
				spriteID = 4
			} else {
				spriteID = 7
			}
		}
	} else {
		seg := ROLL_COOLDOWN / 5
		var rollDir int = -1
		if ship.vel.X > 0 {
			rollDir = 1
		}

		spriteID = 4 + uint16(currentRollCooldown/seg*uint16(rollDir)&3)
	}

	drawSprite(&ship.sprite, ship.pos, 0, spriteID)
}

func unloadTrailPart(ID int) {
	shipTrail[ID] = shipTrail[len(shipTrail)-1]
	shipTrail = shipTrail[:len(shipTrail)-1]
}

func updateShipTrail(shipPosX float64) {
	if len(shipTrail) == 0 {
		return
	}

	for i := 0; i < len(shipTrail); i++ {
		shipTrail[i].mask.A -= 20
		if shipTrail[i].mask.A < 20 || math.Abs(shipTrail[i].pos.X-shipPosX) > 6 {
			unloadTrailPart(i)
			continue
		}

		shipTrail[i].pos = shipTrail[i].pos.Sub(pixel.V(0, (globalVelocity - 15)))
		scale := 2 * (float64(shipTrail[i].mask.A) / 255)

		// Draw trail
		pixel.NewSprite(trailSheet, trailSheet.Bounds()).Draw(
			trailBatch, pixel.IM.Scaled(pixel.ZV, scale).Moved(shipTrail[i].pos))
	}

	trailBatch.Draw(win)
	trailBatch.Clear()
}

func handleRolling() {
	sign := signbit(ship.vel.X)
	if currentRollCooldown == 0 {
		if input.roll && math.Abs(ship.vel.X) > 0.5 {
			ship.vel.X += 8 * sign
			currentRollCooldown = ROLL_COOLDOWN
		}
	} else {
		input.dir.X = 0
		currentRollCooldown--
	}
}
