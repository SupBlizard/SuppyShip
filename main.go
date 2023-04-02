package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

// Globals
const BOUNDARY_STRENGTH float64 = 2
const AXIS_DEADZONE float64 = 0.1
const DEFAULT_GLOBAL_VELOCITY float64 = 10
const ROLL_COOLDOWN int = 20

var WINSIZE pixel.Vec = pixel.V(512, 768)

// Top Bottom Right Left
var borderRanges = [3]float64{400, 150, 70}
var zeroBorder = [3]float64{0, 0, 0}
var frameCount int

var rollCooldown int
var globalVelocity float64 = 5
var globalAcceleration = [2]float64{1.4, 0.6}
var currentLevel uint8

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
	icon := loadPicture("assets/icon.png")
	var iconArr = []pixel.Picture{icon}
	var cfg = pixelgl.WindowConfig{
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

	// Load text atlas
	var textAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

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

	// Load projectile sprite positions
	loadProjectileSpritePos()

	// Generate the star background
	loadStarPhases()
	generateStars()

	var mainColor = color.RGBA{89, 232, 248, 255}
	titleText := text.New(pixel.V(50, WINSIZE.Y-100), textAtlas)
	titleText.Color = mainColor
	fmt.Fprintln(titleText, "Suppy Ship")

	pauseText := text.New(pixel.V(50, WINSIZE.Y-50), textAtlas)
	pauseText.Color = mainColor
	fmt.Fprintln(pauseText, "Paused")

	var (
		paused bool
	)

	for !win.Closed() {
		if currentLevel == 0 {
			titleText.Draw(win, pixel.IM.Scaled(titleText.Orig, 4))
			if win.Pressed(pixelgl.KeyEnter) || win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonStart) {
				currentLevel = 1
			}
		} else if win.JustPressed(pixelgl.KeyEscape) || win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonStart) {
			// Handle pause button
			paused = !paused
			pauseText.Draw(win, pixel.IM.Scaled(pauseText.Orig, 2))
		}

		// Game handling
		if !paused && currentLevel != 0 {
			win.Clear(color.RGBA{0, 0, 0, 0})
			globalVelocity = DEFAULT_GLOBAL_VELOCITY

			// Handle input
			inputDirection, shooting, rolling := handleInput(win)

			// Rolling
			sign := signbit(ship.phys.vel.X)
			if rollCooldown == 0 {
				if rolling && math.Abs(ship.phys.vel.X) > 0.5 {
					ship.phys.vel.X += 9 * sign
					rollCooldown = ROLL_COOLDOWN
				}
			} else {
				inputDirection.X = 0
				ship.phys.vel.X += 0.3 * sign
				rollCooldown--
			}

			// Update ship
			ship.phys = updateShipPhys(ship.phys, inputDirection, rolling)

			// Change ship direction sprite
			ship.sprite.current = 0
			if math.Abs(inputDirection.X) > AXIS_DEADZONE {
				if ship.phys.vel.X > 0 {
					ship.sprite.current = 2
				} else {
					ship.sprite.current = 1
				}
			}

			// Fire bullets
			if shooting && gunCooldown == 0 && skipFrames(reloadDelay) {
				fireBullet(ship.phys.pos)
			}

			// Draw stars
			updateStars()

			// Update Projectiles
			updateProjectiles()

			// Draw ship
			drawSprite(&ship.sprite, ship.phys.pos)

			frameCount++
		}

		// Update window
		win.Update()

	}
}

// Update the ship velocity and position
func updateShipPhys(ship physObj, inputDirection pixel.Vec, rollButton bool) physObj {
	// Give velocity a minimum limit or apply friction to the velocity if there is any
	if ship.vel.Len() <= 0.01 {
		ship.vel = pixel.ZV
	} else {
		ship.vel = ship.vel.Scaled(ship.frc)
	}

	// Add new velocity if there is input
	if inputDirection != pixel.ZV {
		ship.vel = ship.vel.Add(inputDirection)
	}

	// Enforce soft boundary on ship
	if borderCollisions := inBounds(ship.pos, borderRanges); borderCollisions != pixel.ZV {
		var borderDepth float64
		var globalAccIdx int

		if borderCollisions.Y == -1 {
			borderDepth = findBorderDepth(WINSIZE.Y-ship.pos.Y, borderRanges[0])
			globalAccIdx = 0
		} else if borderCollisions.Y == 1 {
			borderDepth = findBorderDepth(ship.pos.Y, borderRanges[1])
			globalAccIdx = 1
		}

		counterAcceleration := ship.acc * BOUNDARY_STRENGTH
		globalVelocity -= (borderDepth * borderCollisions.Y * (DEFAULT_GLOBAL_VELOCITY * globalAcceleration[globalAccIdx]))
		ship.vel.Y += counterAcceleration * borderDepth * borderCollisions.Y

		if borderCollisions.X == -1 {
			ship.vel.X -= counterAcceleration * findBorderDepth(WINSIZE.X-ship.pos.X, borderRanges[2])
		} else if borderCollisions.X == 1 {
			ship.vel.X += counterAcceleration * findBorderDepth(ship.pos.X, borderRanges[2])
		}
	}

	// Add new velocity to the position
	if ship.vel.Len() != 0 {
		ship.pos = ship.pos.Add(ship.vel)
	}

	return ship
}

// Calculate how far into the border something is
func findBorderDepth(pos float64, borderRange float64) float64 {
	return 1 - pos/borderRange
}

// Check if pos is in bounds
func inBounds(pos pixel.Vec, boundaryRange [3]float64) pixel.Vec {
	var boundCollision pixel.Vec = pixel.ZV
	if pos.Y >= WINSIZE.Y-boundaryRange[0] {
		boundCollision.Y = -1
	} else if pos.Y <= boundaryRange[1] {
		boundCollision.Y = 1
	}
	if pos.X >= WINSIZE.X-boundaryRange[2] {
		boundCollision.X = -1
	} else if pos.X <= boundaryRange[2] {
		boundCollision.X = 1
	}

	return boundCollision
}

// Handle user input for a single frame
func handleInput(win *pixelgl.Window) (pixel.Vec, bool, bool) {
	// Initialize dirVec with joystick pos
	var dirVec pixel.Vec = pixel.Vec{
		X: win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftX),
		Y: win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftY) * -1,
	}

	// Add keyboard input to the direction vector
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

	// Trim directional vector and add deadzone
	if dirVecLen := dirVec.Len(); dirVecLen > 1 {
		dirVec = dirVec.Scaled(1 / dirVecLen)
	} else if dirVecLen < AXIS_DEADZONE {
		dirVec = pixel.ZV
	}

	// shoot
	var shootButton bool
	if win.JoystickPressed(pixelgl.Joystick1, pixelgl.ButtonA) || win.Pressed(pixelgl.KeyP) {
		shootButton = true
	}

	// Ignore X axis input if rolling
	var rollButton bool
	if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonX) || win.JustPressed(pixelgl.KeyO) {
		rollButton = true
	}

	return dirVec, shootButton, rollButton
}

func signbit(x float64) float64 {
	sign := 1.0
	if math.Signbit(x) {
		sign = -1.0
	}
	return sign
}

func skipFrames(skip int) bool {
	return frameCount%skip == 0
}

// Lonely Main Function :(
func main() { pixelgl.Run(run) }
