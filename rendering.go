package main

import (
	"image"
	_ "image/png"
	"math/rand"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const GLOBAL_ALTER_SPEED int = 5

var win *pixelgl.Window = nil

type spriteSheet struct {
	scale       float64
	bounds      pixel.Rect
	alterOffset int
	isAltered   bool
	alterSpeed  int
	current     int
	sheet       []*pixel.Sprite
}

const STAR_AMOUNT int = 32
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
	for i := 0; i < STAR_AMOUNT; i++ {

		starArray[i] = star{
			pos:   pixel.V(float64(rand.Int()%int(WINSIZE.X)), float64(rand.Int()%int(WINSIZE.Y))),
			phase: rand.Int() % 6,
			shine: 1,
		}

		if rand.Int()%2 == 0 {
			starArray[i].shine = -1
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
	if frameCount%sprite.alterSpeed == 0 {
		sprite.isAltered = !sprite.isAltered
	}
	if sprite.isAltered {
		sprite.current += sprite.alterOffset
	}

	sprite.sheet[sprite.current].Draw(win, pixel.IM.Scaled(pixel.ZV, sprite.scale).Moved(pos))
}

// Load a spritesheet and return it
func loadSpritesheet(imagePath string, spriteSize pixel.Vec, scale float64) spriteSheet {
	image := loadPicture(imagePath)

	sprite := spriteSheet{
		scale:       scale,
		bounds:      image.Bounds(),
		alterOffset: 0,
		isAltered:   false,
		alterSpeed:  GLOBAL_ALTER_SPEED,
		current:     0,
		sheet:       nil,
	}

	counter := 0
	for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y += spriteSize.Y {
		for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x += spriteSize.X {
			sprite.sheet = append(sprite.sheet, pixel.NewSprite(image, pixel.R(x, y, x+spriteSize.X, y+spriteSize.Y)))
			counter++
		}
	}

	sprite.alterOffset = counter / 2
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
