package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const BULLET_ALLOC_SIZE int = 256
const ONYX_CLUSTER_REQUIREMENT int = 6
const ONYX_CLUSTER_RADIUS float64 = 50

// Projectile Allocation Array
var projectiles [BULLET_ALLOC_SIZE]projectile


// Structs
type projectile struct {
    name string
    phys physObj
    loaded bool
    friendly bool
    scale float64
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
        scale: 4,
    },
    projectile {
        name: "Onyx Bullet",
        phys: shipBulletPhys,
        loaded: true,
        friendly: true,
        scale: 20,
    },
}


// [Create a bullet if a slot is free]
func createBullet(bullets *[BULLET_ALLOC_SIZE]projectile, shipPos pixel.Vec) {
    // Loop through bullet array
    for i:=0;i<BULLET_ALLOC_SIZE;i++ {
        if bullets[i].loaded == false {
            bullets[i] = projectileTypes[0]
            bullets[i].phys.pos = bullets[i].phys.pos.Add(shipPos)
            
            indicies, count := bulletsWithinRadius(bullets, bullets[i].phys.pos, ONYX_CLUSTER_RADIUS)
            if count >= ONYX_CLUSTER_REQUIREMENT {
                freeProjectiles(indicies)
                bulletPos := bullets[i].phys.pos
                bullets[i] = projectileTypes[1]
                bullets[i].phys.pos = bulletPos
            }
            return
        }
    }
}

func freeProjectiles(bull []int) {
    for i:=0;i<len(bull);i++ {
        projectiles[bull[i]].loaded = false
    }
}

func bulletsWithinRadius(bullets *[BULLET_ALLOC_SIZE]projectile, point pixel.Vec, radius float64) ([]int, int) {
    var insideRadius []int
    var bulletCount int = 0
    for i:=0;i<BULLET_ALLOC_SIZE;i++ {
        if bullets[i].loaded && bullets[i].phys.pos.Sub(point).Len() < radius {
            insideRadius = append(insideRadius, i)
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
            bullets[i].sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, bullets[i].scale).Moved(bullets[i].phys.pos))
        } 
    }
}