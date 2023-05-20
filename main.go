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

// var ms uint32
var state uint8
var textAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

// Main
func run() {
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

	// Load things
	loadFragmentSprites()
	loadStarPhases()
	loadStarFields()

	// temp loading enemies (remove when levels are added)
	loadEnemy(0, pixel.V(0, WINY), pixel.V(1, -1))
	loadEnemy(1, pixel.V(200, 400), pixel.ZV)

	// temp loading text fields
	var mainColor = color.RGBA{0, 255, 152, 255}
	titleText := text.New(pixel.V(50, WINY-100), textAtlas)
	pauseText := text.New(pixel.V(50, WINY-50), textAtlas)
	powerText := text.New(pixel.V(50, 50), textAtlas)
	fmt.Fprintln(titleText, TITLE)
	fmt.Fprintln(pauseText, "Paused")
	pauseText.Color = mainColor
	powerText.Color = mainColor
	titleText.Color = mainColor

	var (
		paused  bool
		fps     uint16
		secChan = time.Tick(time.Second)
	)

	for !win.Closed() {
		win.Clear(color.RGBA{0, 0, 0, 0})

		switch state {
		case 0: // Start menu
			if startButton() {
				state = 1
				win.SetCursorVisible(false)
			} else {
				startScreen(titleText)
			}
		case 1: // Playing the game
			if pauseButton() {
				paused = !paused
			}
			if !paused {
				mainGame()
			} else {
				pauseMenu(pauseText)
			}
		}

		// Display framerate
		fps++
		select {
		case <-secChan:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", TITLE, fps))
			fps = 0
		default:
		}

		// Update window
		win.Update()
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
		if input.shoot && skipFrames(ship.reload) && !ship.recharge && ship.heat == 0 {
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

// Handle start screen
func startScreen(titleText *text.Text) {

	// Draw Title
	titleText.Draw(win, pixel.IM.Scaled(titleText.Orig, 4))

	// Draw stars
	updateStars(0, 5, color.RGBA{0, 255, 152, 255})
	frameCount++
}

// Handle pause menu
func pauseMenu(pauseText *text.Text) {
	pauseText.Draw(win, pixel.IM.Scaled(pauseText.Orig, 2))
}

// Lonely Main Function :( even suppy ignores it ):
func main() { pixelgl.Run(run) }
