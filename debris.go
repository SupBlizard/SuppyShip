package main

import (
	"github.com/faiface/pixel"
)

var debrisSprSize pixel.Vec = pixel.V(6, 10)
var debrisSheet pixel.Picture = loadPicture("assets/fragment-spritesheet.png")
var debrisSpritePos = [3][3]pixel.Rect{
	[3]pixel.Rect(batchSpritePos(0, debrisSheet, debrisSprSize)),
	[3]pixel.Rect(batchSpritePos(1, debrisSheet, debrisSprSize)),
}

// Load a piece of debris if there is space
func loadDebris(newDebris debris) {
	if debrislen := uint16(len(debrisAlloc)); debrislen < DEBRIS_ALLOC_SIZE {
		debrisAlloc = append(debrisAlloc, newDebris)
	}
}

// Unload debris
func unloadDebris(idx uint16) {
	debrisAlloc[idx] = debrisAlloc[len(debrisAlloc)-1]
	debrisAlloc = debrisAlloc[:len(debrisAlloc)-1]
}

func updateDebris() {
	if len(debrisAlloc) == 0 {
		return
	}

	// Loop through loaded indexes
	var lenDebris uint16 = uint16(len(debrisAlloc))
	for i := lenDebris - 1; i < lenDebris; i-- {
		if inBounds(debrisAlloc[i].pos, spawnBorder) != pixel.ZV {
			unloadDebris(i)
			continue
		}

		debrisAlloc[i].pos = debrisAlloc[i].pos.Add(debrisAlloc[i].vel)
		debrisAlloc[i].rot += debrisAlloc[i].rotVel
		if debrisAlloc[i].rot > REVOLUTION {
			debrisAlloc[i].rot -= REVOLUTION
		} else if debrisAlloc[i].rot < 0 {
			debrisAlloc[i].rot += REVOLUTION
		}

		// Draw Piece of Debris
		pixel.NewSprite(debrisSheet, debrisSpritePos[debrisAlloc[i].ID[0]][debrisAlloc[i].ID[1]]).Draw(
			debrisBatch, pixel.IM.Scaled(pixel.ZV, debrisAlloc[i].scale).Rotated(pixel.ZV, debrisAlloc[i].rot).Moved(debrisAlloc[i].pos),
		)
	}

	// Draw debris batch
	debrisBatch.Draw(win)
	debrisBatch.Clear()
}
