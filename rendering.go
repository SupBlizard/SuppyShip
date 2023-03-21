package main

import (
    "image"
	"os"

	_ "image/png"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const GLOBAL_ALTER_SPEED int = 5
var win *pixelgl.Window = nil
var spritesAltered bool = false


type spriteSheet struct {
    scale       float64
    bounds      pixel.Rect
    alterOffset int
    isAltered   bool
    alterSpeed  int
    current     int
    sheet []*pixel.Sprite
}

func drawSprite(sprite *spriteSheet, pos pixel.Vec) {
    if frameCount % sprite.alterSpeed == 0 {
        sprite.isAltered = !sprite.isAltered
    }
    if spritesAltered {sprite.current+=sprite.alterOffset}
    
    sprite.sheet[sprite.current].Draw(win, pixel.IM.Scaled(pixel.ZV, sprite.scale).Moved(pos))
}


func loadSpritesheet(imagePath string, spriteSize pixel.Vec, scale float64) (spriteSheet) {
    image, err := loadPicture(imagePath)
    if err != nil {panic(err)}
    sprite := spriteSheet {
        scale: scale,
        bounds: image.Bounds(),
        alterOffset: 0,
        isAltered: false,
        alterSpeed: GLOBAL_ALTER_SPEED,
        current: 0,
        sheet: nil,
    }
    
    counter := 0
    for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y += spriteSize.Y {
        for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x += spriteSize.X {
			sprite.sheet = append(sprite.sheet, pixel.NewSprite(image, pixel.R(x, y, x+spriteSize.X, y+spriteSize.Y)))
            counter++
        }
	}
	
	sprite.alterOffset = counter/2
	println(counter)
	return sprite
}


// [Load a picture from a path]
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {return nil, err}
	defer file.Close()
	
	img, _, err := image.Decode(file)
	if err != nil {return nil, err}
	
	return pixel.PictureDataFromImage(img), nil
}