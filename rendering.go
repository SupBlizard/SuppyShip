package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
)

// Draw sprites
func drawSprite(sprite *spriteSheet, pos pixel.Vec, id uint16) {
	if uint16(frameCount)%sprite.cycleSpeed == 0 {
		sprite.cycle++
		if sprite.cycle >= sprite.cycleNumber {
			sprite.cycle = 0
		}
	}
	sprite.current = id + sprite.offset*sprite.cycle

	sprite.sheet[sprite.current].Draw(win, pixel.IM.Scaled(pixel.ZV, sprite.scale).Moved(pos))
}

// Load a spritesheet and return it
func loadSpritesheet(imagePath string, spriteSize pixel.Vec, scale float64, cycleSpeed uint16) spriteSheet {
	sheet := loadPicture(imagePath)
	sprite := spriteSheet{
		scale:       scale,
		offset:      0,
		cycle:       0,
		cycleNumber: 0,
		cycleSpeed:  cycleSpeed,
		current:     0,
		sheet:       nil,
	}

	var spriteNumber uint16
	for y := sheet.Bounds().Min.Y; y < sheet.Bounds().Max.Y; y += spriteSize.Y {
		for x := sheet.Bounds().Min.X; x < sheet.Bounds().Max.X; x += spriteSize.X {
			sprite.sheet = append(sprite.sheet, pixel.NewSprite(sheet, pixel.R(x, y, x+spriteSize.X, y+spriteSize.Y)))
			spriteNumber++
		}
		sprite.cycleNumber++
	}

	sprite.offset = spriteNumber / sprite.cycleNumber
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
