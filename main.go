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

var (
	gameClock int64 = time.Now().UnixMilli()
	lastClock int64 = gameClock
	dt        float64
)

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
	loadEnemy(0, pixel.V(0, WINY), pixel.V(60, -60))
	loadEnemy(1, pixel.V(200, 400), pixel.ZV)

	// temp loading text fields
	var mainColor = color.RGBA{0, 255, 152, 255}
	titleText := text.New(pixel.V(50, WINY-100), textAtlas)
	pauseText := text.New(pixel.V(50, WINY-50), textAtlas)
	powerText := text.New(pixel.V(50, 50), textAtlas)
	pauseText.Color = mainColor
	powerText.Color = mainColor
	titleText.Color = mainColor

	fmt.Fprintln(titleText, TITLE)
	fmt.Fprintln(pauseText, "Paused")

	var (
		fps    uint16
		paused bool
	)

	for !win.Closed() {
		clockTick()
		// fmt.Printf("%d %d %f\n", gameClock, gameClock-lastClock, dt)

		switch state {
		case 0: // Start menu
			if startButton() {
				state = 1
				win.SetCursorVisible(false)
			} else {
				win.Clear(color.RGBA{0, 0, 0, 0})
				startScreen(titleText)
			}
		case 1: // Playing the game
			if pauseButton() {
				paused = !paused
			}
			if !paused {
				win.Clear(color.RGBA{0, 0, 0, 0})
				mainGame()
			} else {
				pauseMenu(pauseText)

				// Retain the time before pausing
				gameClock = lastClock
			}
		}

		// Calculate framerate
		fps++
		if timePassed(1000) {
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", TITLE, fps))
			fps = 0
		}

		// Update window
		win.Update()
	}
}

func clockTick() {
	lastClock = gameClock
	gameClock = time.Now().UnixMilli()
	dt = float64(gameClock-lastClock) / 1000
}

func mainGame() {
	globalVelocity = 0

	if ship.alive {
		// Update frame's input
		handleInput(win)

		// Update ship
		updateShip()

		// Fire bullets
		if input.shoot && timePassed(ship.reload) && !ship.recharge && ship.heat == 0 {
			if ship.power > 5 {
				fireBullet(ship.pos)
				ship.power -= 5
			} else {
				ship.recharge = true
			}
		}

		// Shield invisibillity frames
		if timePassed(33) && ship.shield.prot > 0 {
			ship.shield.prot--
		}
	}

	// Draw stars
	updateStars(33, 64, nil)

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
	if ship.power < 0xFF && timePassed(33) {
		ship.power++
	}
	if ship.recharge && ship.power > 30 {
		ship.recharge = false
	}
}

// Handle start screen
func startScreen(titleText *text.Text) {

	// Draw Title
	titleText.Draw(win, pixel.IM.Scaled(titleText.Orig, 4))

	// Draw stars
	updateStars(0, 64, color.RGBA{0, 255, 152, 255})
}

// Handle pause menu
func pauseMenu(pauseText *text.Text) {
	pauseText.Draw(win, pixel.IM.Scaled(pauseText.Orig, 2))
}

// Lonely Main Function :( even suppy ignores it ):
func main() { pixelgl.Run(run) }
