package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const BULLET_ALLOC_SIZE int = 256
const ONYX_CLUSTER_REQUIREMENT int = 7
const ONYX_CLUSTER_RADIUS float64 = 30
const ONYX_COOLDOWN int = 60
var reloadDelay int = 4
var gunCooldown = 0

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

var projectileTypes = []projectile {
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
func createBullet(shipPos pixel.Vec) {
    for i:=0;i<BULLET_ALLOC_SIZE;i++ {
        if projectiles[i].loaded == false {
            projectiles[i] = projectileTypes[0]
            projectiles[i].phys.pos = projectiles[i].phys.pos.Add(shipPos)
            
            indicies, count := bulletsWithinRadius(projectiles[i].phys.pos, ONYX_CLUSTER_RADIUS)
            if count >= ONYX_CLUSTER_REQUIREMENT {
                unloadProjectiles(indicies)
                bulletPos := projectiles[i].phys.pos
                projectiles[i] = projectileTypes[1]
                projectiles[i].phys.pos = bulletPos
                
                gunCooldown = ONYX_COOLDOWN
            }
            return
        }
    }
}

// [Unload all selected projectiles]
func unloadProjectiles(selected []int) {
    for i:=0;i<len(selected);i++ {
        projectiles[selected[i]].loaded = false
    }
}

// [Return all of the bullets within a certain radius around a point]
func bulletsWithinRadius(point pixel.Vec, radius float64) ([]int, int) {
    var insideRadius []int
    var projectileCount int = 0
    for i:=0;i<BULLET_ALLOC_SIZE;i++ {
        if projectiles[i].loaded && projectiles[i].phys.pos.Sub(point).Len() < radius {
            insideRadius = append(insideRadius, i)
            projectileCount++
        }
    }
    return insideRadius, projectileCount
}

// [Update the state of each bullet for one frame]
func updateBullets(win *pixelgl.Window) {
    // Update bullets
    for i:=0;i<BULLET_ALLOC_SIZE;i++ {
        if projectiles[i].loaded == false {
            continue
        }
        if inBounds(projectiles[i].phys.pos, WINSIZE, NULL_BOUNDARY_RANGE) != pixel.ZV {
            projectiles[i].loaded = false
        } else {
            projectiles[i].phys.pos = projectiles[i].phys.pos.Add(projectiles[i].phys.vel)
            projectiles[i].sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, projectiles[i].scale).Moved(projectiles[i].phys.pos))
        } 
    }
}