package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

const STAR_MAX_PHASE int = 5
const STAR_DISTANCE float64 = 65
const STAR_RANDOMNESS int = 60
const STAR_SIZE float64 = 5

var starSheet pixel.Picture = loadPicture("assets/star-spritesheet.png")
var starPhases [STAR_MAX_PHASE + 1]pixel.Rect
var starBatch = [2][STAR_MAX_PHASE]*pixel.Batch{
	{pixel.NewBatch(&pixel.TrianglesData{}, starSheet)}, {},
}

// Adjust star distance to fit screen
var starDistance = pixel.Vec{
	X: winsize.X / math.Floor(winsize.X/STAR_DISTANCE),
	Y: winsize.Y / math.Floor(winsize.Y/STAR_DISTANCE),
}

var stars []star = make([]star, 0, 128)

type star struct {
	pos   pixel.Vec
	phase int
	shine int
}

func loadStarPhases() {
	for i := 0; i <= STAR_MAX_PHASE; i++ {
		phase := float64(i) * STAR_SIZE
		starPhases[i] = pixel.R(phase, 0, phase+STAR_SIZE, STAR_SIZE)
	}
}

// Generate stars
func generateStars() {
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
}

func randomVector(limit int) pixel.Vec {
	return pixel.V(float64(rand.Int()%limit), float64(rand.Int()%limit))
}

func updateStarPhases(stars []star) {
	for i := range stars {
		if stars[i].phase >= STAR_MAX_PHASE {
			stars[i].shine = -1
		} else if stars[i].phase < 1 {
			stars[i].shine = 1
		}

		stars[i].phase += 1 * stars[i].shine
	}

}

func updateStars() {
	if skipFrames(2) {
		starBatch[0][0].Clear()

		if skipFrames(4) {
			updateStarPhases(stars)
		}

		for i := 0; i < len(stars); i++ {
			stars[i].pos.Y -= globalVelocity
			if stars[i].pos.Y < 0 {
				stars[i].pos.Y += winsize.Y
			}

			star := pixel.NewSprite(starSheet, starPhases[stars[i].phase])
			star.Draw(starBatch[0][0], pixel.IM.Scaled(pixel.ZV, 2).Moved(stars[i].pos))

		}
	}

	starBatch[0][0].Draw(win)
}
