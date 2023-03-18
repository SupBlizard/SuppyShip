package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const BULLET_ALLOC_SIZE int = 256
const ONYX_CLUSTER_SIZE int = 7

// Projectile Allocation Array
var projectiles [BULLET_ALLOC_SIZE]projectile


// Structs
type projectile struct {
    name string
    phys physObj
    loaded bool
    friendly bool
    sprite *pixel.Sprite
}

var shipBulletPhys physObj = physObj{
    pos: pixel.Vec{0,2},
    vel: pixel.Vec{0,12},
    acc: 0,
    frc: 0,
}

var projectileTypes = [8]projectile {
    projectile {
        name: "Bullet",
        phys: shipBulletPhys,
        loaded: true,
        friendly: true,
    },
    projectile {
        name: "Onyx Bullet",
        phys: shipBulletPhys,
        loaded: true,
        friendly: true,
    },
}


// [Create a bullet if a slot is free]
func createBullet(bullets *[BULLET_ALLOC_SIZE]projectile, shipPos pixel.Vec) {
    // Loop through bullet array
    for i:=0;i<BULLET_ALLOC_SIZE;i++ {
        if bullets[i].loaded == false {
            bullets[i] = projectileTypes[0]
            bullets[i].phys.pos = bullets[i].phys.pos.Add(shipPos)
            
            indicies, count := bulletsWithinRadius(bullets, bullets[i].phys.pos, 30)
            if count >= ONYX_CLUSTER_SIZE {
                // TODO: create freeProjectiles function
                freeProjectiles(indicies)
                bulletPos := bullets[i].phys.pos
                bullets[i] = projectileTypes[1]
                bullets[i].phys.pos = bulletPos
            }
            return
        }
    }
}

func bulletsWithinRadius(bullets *[BULLET_ALLOC_SIZE]projectile, point pixel.Vec, radius float64) ([BULLET_ALLOC_SIZE]int, int) {
    var insideRadius [BULLET_ALLOC_SIZE]int
    var bulletCount int = 0
    for i:=0;i<BULLET_ALLOC_SIZE;i++ {
        if bullets[i].phys.pos.Sub(point).Len() < radius {
            insideRadius[bulletCount] = i
            bulletCount++
        }
    }
    return insideRadius, bulletCount
}

// [Update states of each bullet for one frame]
func updateBullets(bullets *[BULLET_ALLOC_SIZE]projectile, win *pixelgl.Window) {
    // Update bullets
    for i:=0;i<BULLET_ALLOC_SIZE;i++ {
        if bullets[i].loaded == false {
            continue
        }
        if inBounds(bullets[i].phys.pos, WINSIZE, NULL_BOUNDARY_RANGE) != pixel.ZV {
            bullets[i].loaded = false
        } else {
            bullets[i].phys.pos = bullets[i].phys.pos.Add(bullets[i].phys.vel)
            bullets[i].sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 4).Moved(bullets[i].phys.pos))
        } 
    }
}