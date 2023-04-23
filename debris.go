package main

import (
	"github.com/faiface/pixel"
)

var fragSprSize pixel.Vec = pixel.V(6, 10)
var fragSheet pixel.Picture = loadPicture("assets/fragment-spritesheet.png")
var fragSpritePos = [3][3]pixel.Rect{
	[3]pixel.Rect(batchSpritePos(0, fragSheet, fragSprSize)),
	[3]pixel.Rect(batchSpritePos(1, fragSheet, fragSprSize)),
}

// Load a piece of debris if there is space
func loadDebris(newDebris fragment) {
	if debrislen := uint16(len(fragments)); debrislen < DEBRIS_ALLOC_SIZE {
		fragments = append(fragments, newDebris)
	}
}

// Unload debris
func unloadDebris(idx uint16) {
	fragments[idx] = fragments[len(fragments)-1]
	fragments = fragments[:len(fragments)-1]
}

func updateDebris() {
	if len(fragments) == 0 {
		return
	}

	// Loop through loaded indexes
	var lenDebris uint16 = uint16(len(fragments))
	for i := lenDebris - 1; i < lenDebris; i-- {
		if inBounds(fragments[i].pos, spawnBorder) != pixel.ZV {
			unloadDebris(i)
			continue
		}

		fragments[i].pos = fragments[i].pos.Add(fragments[i].vel)
		fragments[i].rot += fragments[i].rotVel
		if fragments[i].rot > REVOLUTION {
			fragments[i].rot -= REVOLUTION
		} else if fragments[i].rot < 0 {
			fragments[i].rot += REVOLUTION
		}

		// Draw Piece of Debris
		pixel.NewSprite(fragSheet, fragSpritePos[fragments[i].ID[0]][fragments[i].ID[1]]).Draw(
			fragmentBatch, pixel.IM.Scaled(pixel.ZV, fragments[i].scale).Rotated(pixel.ZV, fragments[i].rot).Moved(fragments[i].pos),
		)
	}

	// Draw debris batch
	fragmentBatch.Draw(win)
	fragmentBatch.Clear()
}
