package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

const STAR_AMOUNT int = 80
const STAR_MAX_PHASE int = 5

var starSheet pixel.Picture = loadPicture("assets/star-spritesheet.png")
var starBatch *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, starSheet)
var starSize pixel.Vec = pixel.V(5, 5)
var stars [STAR_AMOUNT]star
var starPhases [STAR_MAX_PHASE + 1]pixel.Rect

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
	var floatStarAmount float64 = float64(STAR_AMOUNT)
	var sqrtStarAmount float64 = math.Sqrt(floatStarAmount)
	var windowRatio float64 = WINSIZE.Y / WINSIZE.X

	// Get Grid size
	closestFactor := 1.0
	closestRatio := floatStarAmount
	for i := 1.0; floatStarAmount/i > sqrtStarAmount; i *= 2 {
		currentRatio := floatStarAmount / math.Pow(i, 2)
		if math.Abs(currentRatio-windowRatio) < closestRatio {
			closestRatio = currentRatio
			closestFactor = i
		}
	}

	// TODO: Need to calculate the optimal scalar (currently hardcoded)
	starGridRatio := pixel.V(closestFactor, floatStarAmount/closestFactor).Scaled(9)

	// Generate stars
	currentStar := 0
	for y := 0.0; y < WINSIZE.Y; y += starGridRatio.Y {
		for x := 0.0; x < WINSIZE.X; x += starGridRatio.X {
			if currentStar >= STAR_AMOUNT {
				return
			}

			stars[currentStar] = star{
				pos:   pixel.V(x, y).Add(pixel.V(float64(rand.Int()%80), float64(rand.Int()%80))),
				phase: rand.Int() % 6,
				shine: 1,
			}

			currentStar++
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
		for i := 0; i < STAR_AMOUNT; i++ {

			if skipFrames(4) {
				updateStarPhase(i)
			}

			stars[i].pos.Y -= globalVelocity
			if stars[i].pos.Y < 0 {
				stars[i].pos.Y += WINSIZE.Y
			}

			star := pixel.NewSprite(starSheet, starPhases[stars[i].phase])
			star.Draw(starBatch, pixel.IM.Scaled(pixel.ZV, 2).Moved(stars[i].pos))

		}
	}

	starBatch.Draw(win)
}
