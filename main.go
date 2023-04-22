package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

// Core globals
var win *pixelgl.Window = nil
var frameCount uint32

// Main
func run() {
	// Self explanatory
	loadStuff()

	// Load text atlas
	var textAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

	var mainColor = color.RGBA{89, 232, 248, 255}
	titleText := text.New(pixel.V(50, WINY-100), textAtlas)
	titleText.Color = mainColor
	fmt.Fprintln(titleText, "Suppy Ship")

	pauseText := text.New(pixel.V(50, WINY-50), textAtlas)
	pauseText.Color = mainColor
	fmt.Fprintln(pauseText, "Paused")

	powerText := text.New(pixel.V(50, 50), textAtlas)
	powerText.Color = mainColor

	// temp add enemy asteroid for testing
	loadEnemy(0, win.Bounds().Center(), pixel.ZV)

	var (
		paused         bool
		safetyRecharge bool

		frames int
		second = time.Tick(time.Second)
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

			// Update frame's input
			handleInput(win)

			// Update ship
			updateShip()

			// Fire bullets
			if input.shoot && skipFrames(reloadDelay) && !safetyRecharge && gunCooldown == 0 {
				if ship.power > 5 {
					fireBullet(ship.pos)
					ship.power -= 5
				} else {
					safetyRecharge = true
				}
			}

			// Draw stars
			updateStars()

			// Update Projectiles
			updateProjectiles()

			// Update Enemies
			updateEnemies()

			// Draw ship
			drawShip()

			// Draw ship trail
			updateShipTrail(ship.pos.X)

			// Increment Ship power
			if ship.power < 0xFF && skipFrames(2) {
				ship.power++
			}
			if safetyRecharge && ship.power > 30 {
				safetyRecharge = false
			}

			fmt.Fprintln(powerText, ship.power)
			powerText.Draw(win, pixel.IM.Scaled(powerText.Orig, 2))
			powerText.Clear()

			frameCount++
		}

		// Update window
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", TITLE, frames))
			frames = 0
		default:
		}

	}
}

func loadStuff() {
	// Create new window
	windowPointer, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  TITLE,
		Bounds: pixel.R(0, 0, WINX, WINY),
		Icon:   []pixel.Picture{loadPicture("assets/icon.png")},
		VSync:  true,
	})

	win = windowPointer
	if err != nil {
		panic(err)
	}

	// Load projectile sprite positions
	projSpritePos := loadSpritePos(projectileSheet, pixel.V(6, 16))
	for y := 0; y < len(projSpritePos); y++ {
		copy(projectileTypes[y].sprite.pos[:], projSpritePos[y])
	}

	// Generate the star background
	loadStarPhases()
	loadStarFields()
}

// Handle user input for a single frame
func handleInput(win *pixelgl.Window) {
	// Initialize dirVec with joystick pos
	input.dir = pixel.Vec{
		X: win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftX),
		Y: win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftY) * -1,
	}

	// Add keyboard input to the direction vector
	if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
		input.dir = input.dir.Add(inputLookup[0])
	}
	if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
		input.dir = input.dir.Add(inputLookup[1])
	}
	if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
		input.dir = input.dir.Add(inputLookup[2])
	}
	if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
		input.dir = input.dir.Add(inputLookup[3])
	}

	// Trim directional vector and add deadzone
	if dirVecLen := input.dir.Len(); dirVecLen > 1 {
		input.dir = input.dir.Scaled(1 / dirVecLen)
	} else if dirVecLen < AXIS_DEADZONE {
		input.dir = pixel.ZV
	}

	// Shoot
	if win.JoystickPressed(pixelgl.Joystick1, pixelgl.ButtonRightBumper) || win.Pressed(pixelgl.KeyP) {
		input.shoot = true
	} else {
		input.shoot = false
	}

	// Roll
	if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonA) || win.JustPressed(pixelgl.KeyO) {
		input.roll = true
	} else {
		input.roll = false
	}
}

// Lonely Main Function :( even suppy ignores it ):
func main() { pixelgl.Run(run) }
