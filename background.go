package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Adjust star distance to fit screen
var starDistance = pixel.Vec{
	X: winsize.X / math.Floor(winsize.X/STAR_DISTANCE),
	Y: winsize.Y / math.Floor(winsize.Y/STAR_DISTANCE),
}

func loadStarPhases() {
	for i := uint8(0); i <= STAR_MAX_PHASE; i++ {
		phase := float64(i) * STAR_SIZE
		starSprites[i] = pixel.NewSprite(starSheet, pixel.R(phase, 0, phase+STAR_SIZE, STAR_SIZE))
	}
}

func loadStarFields() {
	var stars = [STARFIELD_NUMBER][]star{generateStars(), generateStars()}

	for i := uint8(0); i < STARFIELD_NUMBER; i++ {
		for j := uint8(0); j < STAR_PHASES; j++ {

			starFields[i][j] = renderStars(stars[i])
			stars[i] = updateStarPhases(stars[i])
		}
	}
}

// Generate stars
func generateStars() []star {
	var stars []star = make([]star, 0, 128)
	for y := 0.0; math.Round(y) < winsize.Y; y += starDistance.Y {
		for x := 0.0; math.Round(x) <= winsize.X; x += starDistance.X {
			// Ignore positions outside the bounds
			randomPos := pixel.V(x, y).Add(randomVector(STAR_RANDOMNESS))
			if inBounds(randomPos, windowBorder) != pixel.ZV {
				continue
			}

			stars = append(stars, star{
				pos:   randomPos,
				phase: rand.Int() % 6,
				shine: 1,
			})
		}
	}
	return stars
}

func renderStars(stars []star) *pixel.Sprite {
	starfield := pixelgl.NewCanvas(pixel.R(0, 0, winsize.X, winsize.Y))
	for _, star := range stars {
		starSprites[star.phase].Draw(starfield, pixel.IM.Scaled(pixel.ZV, 2).Moved(star.pos))
	}

	return pixel.NewSprite(starfield, starfield.Bounds())
}

func updateStarPhases(stars []star) []star {
	for i := range stars {
		if stars[i].phase >= int(STAR_MAX_PHASE) {
			stars[i].shine = -1
		} else if stars[i].phase < 1 {
			stars[i].shine = 1
		}

		stars[i].phase += stars[i].shine
	}
	return stars
}

func updateStars() {
	for i := uint8(0); i < STARFIELD_NUMBER; i++ {
		currentPhase := (frameCount / 4) % int(STAR_PHASES)
		if skipFrames(2) {
			starfieldPos[i] = starfieldPos[i].Sub(pixel.V(0, globalVelocity))
			if starfieldPos[i].Y < winsize.Y*-0.5 {
				starfieldPos[i].Y = winsize.Y * 1.5
			}
		}

		starFields[i][currentPhase].Draw(win, pixel.IM.Moved(starfieldPos[i]))
	}
}
