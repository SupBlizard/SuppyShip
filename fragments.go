package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

var fragSprSize pixel.Vec = pixel.V(6, 10)
var fragSheet pixel.Picture = loadPicture("assets/fragment-spritesheet.png")
var fragSpritePos = [3][3]pixel.Rect{}

func loadFragmentSprites() {
	max := uint8((fragSheet.Bounds().Max.X - fragSheet.Bounds().Min.X) / fragSprSize.X)
	for i := uint8(0); i < max; i++ {
		fragSpritePos[i] = [3]pixel.Rect(batchSpritePos(i, fragSheet, fragSprSize))
	}
}

// Spawn a cluster of fragments
func fragmentObject(info *fragInfo, seg []uint8, pos pixel.Vec, vel pixel.Vec) {
	vectors, angles := spreadFragments(info.frags)

	for i, vec := range vectors {
		if len(seg) == i {
			seg = append(seg, uint8(rand.Int31()%3))
		}
		loadFragment(fragment{
			ID:     [2]uint8{info.ID, seg[i]},
			pos:    pos.Add(vec.Scaled(info.radius)),
			vel:    vec.Add(vel).Scaled(info.power),
			rot:    angles[i],
			rotVel: rand.Float64() - 0.5,
			scale:  info.scale,
		})
	}
}

// Spread out fragment directions evenly
func spreadFragments(n uint8) ([]pixel.Vec, []float64) {
	var points []pixel.Vec
	var angles []float64
	var spread float64 = REVOLUTION / float64(n)
	for i := 0.0; uint8(i) < n; i++ {
		points = append(points, pixel.V(math.Cos(spread*i), math.Sin(spread*i)))
		angles = append(angles, spread*i)
	}

	return points, angles
}

// Load a fragment if there is space
func loadFragment(newDebris fragment) {
	if debrislen := uint16(len(fragments)); debrislen < DEBRIS_ALLOC_SIZE {
		fragments = append(fragments, newDebris)
	}
}

// Unload fragment
func unloadFragment(idx uint16) {
	fragments[idx] = fragments[len(fragments)-1]
	fragments = fragments[:len(fragments)-1]
}

func updateFragments() {
	if len(fragments) == 0 {
		return
	}

	// Loop through loaded indexes
	var lenDebris uint16 = uint16(len(fragments))
	for i := lenDebris - 1; i < lenDebris; i-- {
		if inBounds(fragments[i].pos, spawnBorder) != pixel.ZV {
			unloadFragment(i)
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
