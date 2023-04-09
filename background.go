package main

import (
	"math/rand"

	"github.com/faiface/pixel"
)

const STAR_MAX_PHASE int = 5

var starSheet pixel.Picture = loadPicture("assets/star-spritesheet.png")
var starBatch *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, starSheet)
var starSize pixel.Vec = pixel.V(5, 5)
var starPhases [STAR_MAX_PHASE + 1]pixel.Rect
var starDistance pixel.Vec = pixel.V(90, 90)
var stars []star = make([]star, 0, 128)

type star struct {
	pos   pixel.Vec
	phase int
	shine int
}

func loadStarPhases() {
	for i := 0; i <= STAR_MAX_PHASE; i++ {
		phase := float64(i) * starSize.X
		starPhases[i] = pixel.R(phase, 0, phase+starSize.X, starSize.Y)
	}
}

func generateStars() {
	var starNumbers pixel.Vec = pixel.V(winsize.X/starDistance.X, winsize.Y/starDistance.Y)
	starDistance = starDistance.Add(pixel.Vec{
		X: starNumbers.X - float64(uint8(starNumbers.X)),
		Y: starNumbers.Y - float64(uint8(starNumbers.Y)),
	})

	// Generate stars
	var shiftRow bool
	var shiftAmount float64 = starDistance.X / 2
	var renderBounds pixel.Vec = pixel.V(starNumbers.X*starDistance.X, starNumbers.Y*starDistance.Y)
	for y := 0.0; y < renderBounds.Y; y += starDistance.Y {
		x := 0.0
		if shiftRow {
			x = shiftAmount
		}

		shiftRow = !shiftRow
		for ; x < renderBounds.X; x += starDistance.X {

			//
			stars = append(stars, star{
				pos:   pixel.V(x, y).Add(pixel.V(float64(rand.Int()%50), float64(rand.Int()%50))),
				phase: rand.Int() % 6,
				shine: 1,
			})
		}
	}
}

func updateStarPhase(star int) {
	if stars[star].phase >= STAR_MAX_PHASE {
		stars[star].shine = -1
	} else if stars[star].phase < 1 {
		stars[star].shine = 1
	}

	stars[star].phase += 1 * stars[star].shine
}

func updateStars() {
	if skipFrames(2) {
		starBatch.Clear()
		for i := 0; i < len(stars); i++ {

			if skipFrames(4) {
				updateStarPhase(i)
			}

			stars[i].pos.Y -= globalVelocity
			if stars[i].pos.Y < 0 {
				stars[i].pos.Y += winsize.Y
			}

			star := pixel.NewSprite(starSheet, starPhases[stars[i].phase])
			star.Draw(starBatch, pixel.IM.Scaled(pixel.ZV, 2).Moved(stars[i].pos))

		}
	}

	starBatch.Draw(win)
}
