package main

import (
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const GLOBAL_ALTER_SPEED int = 5

var win *pixelgl.Window = nil

type spriteSheet struct {
	offset      int
	cycle       int
	cycleNumber int
	cycleSpeed  int
	current     int
	scale       float64
	sheet       []*pixel.Sprite
}

const STAR_AMOUNT int = 80
const STAR_MAX_PHASE int = 5

type star struct {
	pos   pixel.Vec
	phase int
	shine int
}

// Background stuff
var starSheet pixel.Picture = loadPicture("assets/star-spritesheet.png")
var starBatch *pixel.Batch = pixel.NewBatch(&pixel.TrianglesData{}, starSheet)
var starSize pixel.Vec = pixel.V(5, 5)
var starArray [STAR_AMOUNT]star
var starPhases [STAR_MAX_PHASE + 1]pixel.Rect

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

			starArray[currentStar] = star{
				pos:   pixel.V(x, y).Add(pixel.V(float64(rand.Int()%80), float64(rand.Int()%80))),
				phase: rand.Int() % 6,
				shine: 1,
			}

			currentStar++
		}
	}
}

func updateStarPhase(star int) {
	if starArray[star].phase >= STAR_MAX_PHASE {
		starArray[star].shine = -1
	} else if starArray[star].phase < 1 {
		starArray[star].shine = 1
	}

	starArray[star].phase += 1 * starArray[star].shine
}

func updateStars() {
	if skipFrames(2) {
		starBatch.Clear()
		for i := 0; i < STAR_AMOUNT; i++ {

			if skipFrames(4) {
				updateStarPhase(i)
			}

			starArray[i].pos.Y -= globalVelocity
			if starArray[i].pos.Y < 0 {
				starArray[i].pos.Y += WINSIZE.Y
			}

			star := pixel.NewSprite(starSheet, starPhases[starArray[i].phase])
			star.Draw(starBatch, pixel.IM.Scaled(pixel.ZV, 2).Moved(starArray[i].pos))

		}
	}

	starBatch.Draw(win)

}

// Draw a sprite
func drawSprite(sprite *spriteSheet, pos pixel.Vec) {
	if frameCount%sprite.cycleSpeed == 0 {
		sprite.cycle++
		if sprite.cycle >= sprite.cycleNumber {
			sprite.cycle = 0
		}
	}
	sprite.current += sprite.offset * sprite.cycle

	sprite.sheet[sprite.current].Draw(win, pixel.IM.Scaled(pixel.ZV, sprite.scale).Moved(pos))
}

// Load a spritesheet and return it
func loadSpritesheet(imagePath string, spriteSize pixel.Vec, scale float64) spriteSheet {
	image := loadPicture(imagePath)

	sprite := spriteSheet{
		scale:       scale,
		offset:      0,
		cycle:       0,
		cycleNumber: 0,
		cycleSpeed:  GLOBAL_ALTER_SPEED,
		current:     0,
		sheet:       nil,
	}

	var spriteNumber int
	var cycleNumber int
	for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y += spriteSize.Y {
		for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x += spriteSize.X {
			sprite.sheet = append(sprite.sheet, pixel.NewSprite(image, pixel.R(x, y, x+spriteSize.X, y+spriteSize.Y)))
			spriteNumber++
		}
		cycleNumber++
	}

	sprite.cycleNumber = cycleNumber
	sprite.offset = spriteNumber / cycleNumber
	return sprite
}

// Load a picture from a path
func loadPicture(path string) pixel.Picture {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	return pixel.PictureDataFromImage(img)
}
