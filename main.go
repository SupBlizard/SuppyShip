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
var frameCount uint16
var textAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

// Main
func run() {
	// Self explanatory
	loadStuff()

	var mainColor = color.RGBA{0, 255, 152, 255}

	titleText := text.New(pixel.V(50, WINY-100), textAtlas)
	pauseText := text.New(pixel.V(50, WINY-50), textAtlas)
	powerText := text.New(pixel.V(50, 50), textAtlas)

	// Write to text
	fmt.Fprintln(titleText, TITLE)
	fmt.Fprintln(pauseText, "Paused")

	// Color text
	pauseText.Color = mainColor
	powerText.Color = mainColor
	titleText.Color = mainColor

	// temp add enemy asteroid for testing
	loadEnemy(0, pixel.V(0, WINY), pixel.V(1, -1))

	// temp add enemy asteroid for testing
	loadEnemy(1, pixel.V(200, 400), pixel.ZV)

	var (
		paused bool
		frames int
		second = time.Tick(time.Second)
	)

	for !win.Closed() {
		win.Clear(color.RGBA{0, 0, 0, 0})

		// Title Screen
		if currentLevel == 0 {
			startScreen(titleText)
		} else if pauseButton() {
			paused = !paused
		}

		if paused {
			pauseMenu(pauseText)
		}

		// Game handling
		if !paused && currentLevel != 0 {
			mainGame()
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

func mainGame() {
	globalVelocity = 0

	if ship.alive {
		// Update frame's input
		handleInput(win)

		// Update ship
		updateShip()

		// Fire bullets
		if input.shoot && skipFrames(reloadDelay) && !ship.recharge && gunCooldown == 0 {
			if ship.power > 5 {
				fireBullet(ship.pos)
				ship.power -= 5
			} else {
				ship.recharge = true
			}
		}

		// Shield invisibillity frames
		if skipFrames(2) && ship.shield.prot > 0 {
			ship.shield.prot--
		}
	}

	// Draw stars
	updateStars(2, 4, nil)

	// Update Projectiles
	updateProjectiles()

	// Update Enemies
	updateEnemies()

	if ship.alive {
		// Draw ship trail
		updateShipTrail(ship.pos)

		// Draw ship
		if ship.shield.prot%5 != 1 {
			drawShip()
		}
		// Draw shield
		if ship.shield.active {
			drawSprite(&ship.shield.sprite, ship.pos, 0, 0)
		}
	}

	// Update Debris
	updateFragments()

	// Increment Ship power
	if ship.power < 0xFF && skipFrames(2) {
		ship.power++
	}
	if ship.recharge && ship.power > 30 {
		ship.recharge = false
	}

	frameCount++
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
	win.SetCursorVisible(false)

	loadFragmentSprites()
	loadStarPhases()
	loadStarFields()
}

// Handle start screen
func startScreen(titleText *text.Text) {
	win.Clear(color.RGBA{0, 0, 0, 0})

	// Draw Title
	titleText.Draw(win, pixel.IM.Scaled(titleText.Orig, 4))

	// Draw stars
	updateStars(0, 5, color.RGBA{0, 255, 152, 255})

	// Start game
	if startButton() {
		currentLevel = 1
	}
	frameCount++
}

// Handle pause menu
func pauseMenu(pauseText *text.Text) {
	pauseText.Draw(win, pixel.IM.Scaled(pauseText.Orig, 2))
}

// Button functions
func pauseButton() bool {
	return win.JustPressed(pixelgl.KeyEscape) || win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonStart)
}

func startButton() bool {
	return win.Pressed(pixelgl.KeyEnter) || win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonStart)
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
	if win.JoystickPressed(pixelgl.Joystick1, pixelgl.ButtonRightBumper) || win.Pressed(pixelgl.MouseButtonLeft) {
		input.shoot = true
	} else {
		input.shoot = false
	}

	// Roll
	if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonA) || win.JustPressed(pixelgl.MouseButtonRight) {
		input.roll = true
	} else {
		input.roll = false
	}
}

// Lonely Main Function :( even suppy ignores it ):
func main() { pixelgl.Run(run) }
