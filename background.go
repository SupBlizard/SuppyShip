package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Adjust star distance to fit screen
var starDistance = pixel.Vec{
	X: WINX / math.Floor(WINX/STAR_DISTANCE),
	Y: WINY / math.Floor(WINY/STAR_DISTANCE),
}

var starFields [STARFIELD_NUMBER][STAR_PHASES]*pixel.Sprite
var starfieldPos = [STARFIELD_NUMBER]pixel.Vec{
	pixel.V(WINX*0.5, WINY*0.5),
	pixel.V(WINX*0.5, WINY*1.5),
}

func loadStarPhases() {
	for i := int8(0); i <= STAR_MAX_PHASE; i++ {
		phase := float64(i) * STAR_SIZE
		starSprites[i] = pixel.NewSprite(starSheet, pixel.R(phase, 0, phase+STAR_SIZE, STAR_SIZE))
	}
}

func loadStarFields() {
	var stars = [STARFIELD_NUMBER][]star{generateStars(), generateStars()}

	for i := uint8(0); i < STARFIELD_NUMBER; i++ {
		for j := int8(0); j < STAR_PHASES; j++ {

			starFields[i][j] = renderStars(stars[i])
			stars[i] = updateStarPhases(stars[i])
		}
	}
}

// Generate stars
func generateStars() []star {
	var stars []star = make([]star, 0, 128)
	for y := 0.0; math.Round(y) < WINY; y += starDistance.Y {
		for x := 0.0; math.Round(x) <= WINX; x += starDistance.X {
			// Ignore positions outside the bounds
			randomPos := pixel.V(x, y).Add(randomVector(STAR_RANDOMNESS))
			if inBounds(randomPos, windowBorder) != pixel.ZV {
				continue
			}

			rnd := int8(rand.Int() % 6)

			stars = append(stars, star{
				pos:   randomPos,
				phase: rnd,
				shine: 1,
			})
		}
	}
	return stars
}

func renderStars(stars []star) *pixel.Sprite {
	starfield := pixelgl.NewCanvas(pixel.R(0, 0, WINX, WINY))
	for _, cStar := range stars {
		starSprites[cStar.phase].Draw(starfield, pixel.IM.Scaled(pixel.ZV, 2).Moved(cStar.pos))
	}

	return pixel.NewSprite(starfield, starfield.Bounds())
}

func updateStarPhases(stars []star) []star {
	for i := range stars {
		if stars[i].phase >= STAR_MAX_PHASE {
			stars[i].shine = -1
		} else if stars[i].phase < 1 {
			stars[i].shine = 1
		}

		stars[i].phase += stars[i].shine
	}
	return stars
}

func updateStars(posRate uint16, phaseRate uint16, colorMask color.Color) {
	for i := uint8(0); i < STARFIELD_NUMBER; i++ {
		// Update starfield pos
		if skipFrames(posRate) {
			starfieldPos[i] = starfieldPos[i].Sub(pixel.V(0, globalVelocity))
			if starfieldPos[i].Y < WINY*-0.5 {
				starfieldPos[i].Y = WINY * 1.5
			}
		}

		// Draw starfield
		starFields[i][(frameCount/phaseRate)%uint16(STAR_PHASES)].DrawColorMask(win, pixel.IM.Moved(starfieldPos[i]), colorMask)
	}
}
