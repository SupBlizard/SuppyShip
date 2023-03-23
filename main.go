package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Globals
const BOUNDARY_STRENGTH float64 = 2
const AXIS_LOWERBOUND float64 = 0.1

var WINSIZE pixel.Vec = pixel.V(500, 700)
var BOUNDARY_RANGE = [4]float64{300, 100, 70, 70} // Top Bottom Right Left
var NULL_BOUNDARY_RANGE = [4]float64{0, 0, 0, 0}
var frameCount int

// Input direction vector lookup table
var inputVec = [4]pixel.Vec{
	pixel.V(0, 1),
	pixel.V(-1, 0),
	pixel.V(0, -1),
	pixel.V(1, 0),
}

// Structs
type physObj struct {
	pos pixel.Vec // position
	vel pixel.Vec // velocity
	acc float64   // acceleration
	frc float64   // friction
}

type player struct {
	phys   physObj
	sprite spriteSheet
}

// Main
func run() {
	icon, err := loadPicture("assets/icon.png")
	if err != nil {
		panic(err)
	}
	var iconArr []pixel.Picture = []pixel.Picture{icon}

	cfg := pixelgl.WindowConfig{
		Title:  "Suppy Ship",
		Bounds: pixel.R(0, 0, WINSIZE.X, WINSIZE.Y),
		VSync:  true,
		Icon:   iconArr,
	}

	// Create new window
	windowPointer, err := pixelgl.NewWindow(cfg)
	win = windowPointer
	if err != nil {
		panic(err)
	}

	// Initialize player ship
	ship := player{
		phys: physObj{
			pos: win.Bounds().Center(),
			vel: pixel.ZV,
			acc: 1.1,
			frc: 1 - 0.08,
		},
		sprite: loadSpritesheet("assets/ship-spritesheet.png", pixel.V(13, 18), 2),
	}

	// Load projectile sprites
	loadProjectileSprites()

	var paused bool = false
	for !win.Closed() {

		// Handle pause button
		if win.JustPressed(pixelgl.KeyEscape) || win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonStart) {
			paused = !paused
		}

		if !paused {
			frameCount++
			win.Clear(colornames.Black)

			// Handle input
			inputDirection, shooting := handleInput(win)

			// Update ship
			ship.phys = updateShipPhys(ship.phys, inputDirection)

			// Create new bullets
			if shooting && gunCooldown == 0 && (frameCount%reloadDelay) == 0 {
				createBullet(ship.phys.pos)
			}
			if gunCooldown > 0 {
				gunCooldown--
			}
			updateBullets()

			// Change ship direction sprite
			ship.sprite.current = 0
			if inputDirection.X != 0 {
				if ship.phys.vel.X > 0 {
					ship.sprite.current = 2
				} else {
					ship.sprite.current = 1
				}
			}

			drawSprite(&ship.sprite, ship.phys.pos)
		}
		// Update window
		win.Update()
	}
}

// Update the ship velocity and position
func updateShipPhys(ship physObj, inputDirection pixel.Vec) physObj {
	// Give velocity a minimum limit or apply friction to the velocity if there is any
	if ship.vel.Len() <= 0.01 {
		ship.vel = pixel.ZV
	} else {
		ship.vel = ship.vel.Scaled(ship.frc)
	}

	// Add new velocity if there is input
	if inputDirection != pixel.ZV {
		//.Scaled(ship.acc / inputDirection.Len())
		ship.vel = ship.vel.Add(inputDirection)
	}

	// Enforce soft boundary on ship
	boundsCollisions := inBounds(ship.pos, WINSIZE, BOUNDARY_RANGE)
	if boundsCollisions != pixel.ZV {
		if boundsCollisions.Y == 1 {
			ship.vel.Y -= borderForce(ship.acc, BOUNDARY_RANGE[0], WINSIZE.Y-ship.pos.Y)
		} else if boundsCollisions.Y == -1 {
			ship.vel.Y += borderForce(ship.acc, BOUNDARY_RANGE[1], ship.pos.Y)
		}
		if boundsCollisions.X == 1 {
			ship.vel.X -= borderForce(ship.acc, BOUNDARY_RANGE[2], WINSIZE.X-ship.pos.X)
		} else if boundsCollisions.X == -1 {
			ship.vel.X += borderForce(ship.acc, BOUNDARY_RANGE[3], ship.pos.X)
		}
	}

	// Add new velocity to the position if there is any input
	if ship.vel.Len() != 0 {
		ship.pos = ship.pos.Add(ship.vel)
	}

	return ship
}

// Border force formula
func borderForce(acc float64, boundaryRange float64, pos float64) float64 {
	return acc * BOUNDARY_STRENGTH * (1 - pos/boundaryRange)
}

// Check if pos is in bounds
func inBounds(pos pixel.Vec, max pixel.Vec, bounds [4]float64) pixel.Vec {
	var boundCollision pixel.Vec = pixel.ZV
	if pos.Y >= max.Y-bounds[0] {
		boundCollision.Y = 1
	} else if pos.Y <= bounds[1] {
		boundCollision.Y = -1
	}
	if pos.X >= max.X-bounds[2] {
		boundCollision.X = 1
	} else if pos.X <= bounds[3] {
		boundCollision.X = -1
	}

	return boundCollision
}

// Handle user input for a single frame
func handleInput(win *pixelgl.Window) (pixel.Vec, bool) {
	var dirVec pixel.Vec = pixel.ZV
	var shootButton bool = false

	if win.JoystickPresent(pixelgl.Joystick1) {
		// Add gamepad axis positions to the direction vector
		dirVec.X = win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftX)
		dirVec.Y = win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftY) * -1

		// Ignore very small values from axes as they could be slight stick drift
		if dirVec.X < AXIS_LOWERBOUND && dirVec.X > AXIS_LOWERBOUND*-1 {
			dirVec.X = 0
		}
		if dirVec.Y < AXIS_LOWERBOUND && dirVec.Y > AXIS_LOWERBOUND*-1 {
			dirVec.Y = 0
		}

		// Shoot
		if win.JoystickPressed(pixelgl.Joystick1, pixelgl.ButtonA) {
			shootButton = true
		}
	} else {
		// Add input vectors to the direction vector
		if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			dirVec = dirVec.Add(inputVec[0])
		}
		if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
			dirVec = dirVec.Add(inputVec[1])
		}
		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			dirVec = dirVec.Add(inputVec[2])
		}
		if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
			dirVec = dirVec.Add(inputVec[3])
		}

		// Shoot
		if win.Pressed(pixelgl.KeySpace) {
			shootButton = true
		}
	}

	return dirVec, shootButton
}

// Lonely Main Function :(
func main() { pixelgl.Run(run) }
