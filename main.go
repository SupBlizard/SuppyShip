package main

import (
	"image"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Globals
var WINSIZE pixel.Vec = pixel.Vec{500,700}
const BOUNDARY_STRENGTH float64 = 2
const BULLET_ALLOC_SIZE int = 256

var BOUNDARY_RANGE =       [4]float64 {300,100,70,70,} // Top Bottom Right Left
var NULL_BOUNDARY_RANGE =  [4]float64 {0,0,0,0,}
var inputVec =              [4]pixel.Vec {
    pixel.Vec{0, 1},
    pixel.Vec{-1,0},
    pixel.Vec{0,-1},
    pixel.Vec{1, 0},
}

// Structs
type physObj struct {
    pos   pixel.Vec     // position
    vel   pixel.Vec     // velocity
    acc   float64       // acceleration
    frc   float64       // friction
}

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

// Main
func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Suppy Ship",
		Bounds: pixel.R(0, 0, WINSIZE.X, WINSIZE.Y),
		VSync:  true,
	}
	
	// Create new window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {panic(err)}
    
    // TODO: make full sprite loader
	shipImage, err := loadPicture("ship.png")
	if err != nil {panic(err)}
    bulletImage, err := loadPicture("bullet.png")
	if err != nil {panic(err)}
	
    var shipSprite = pixel.NewSprite(shipImage, shipImage.Bounds())
    var ship physObj = physObj {
        pos: win.Bounds().Center(),
        vel: pixel.ZV,
        acc: 1.1,
        frc: 1 - 0.08,
    }
     
    var reloadDelay uint8 = 4
    var projectiles [BULLET_ALLOC_SIZE]projectile
    
    // This is temporary
    projectileTypes[0].sprite = pixel.NewSprite(bulletImage, bulletImage.Bounds())
    projectileTypes[1].sprite = pixel.NewSprite(bulletImage, bulletImage.Bounds())


    var frameCount uint32 = 0
    
    var paused bool = false
    var lastPauseButtonState bool = false
    
	for !win.Closed() {
        // Handle keyboard input
        inputDirection, shooting, pauseButton := handleInput(win)
        
        if pauseButton == true && lastPauseButtonState != true {paused = !paused}
        lastPauseButtonState = pauseButton
        
        if !paused {
            frameCount++
            win.Clear(colornames.Black)
            
            // Update ship
            ship = updateShip(ship, inputDirection)
            
            // Create new bullets
            if shooting && (frameCount % uint32(reloadDelay)) == 0 {
                createBullet(&projectiles, ship.pos)
            }
            updateBullets(&projectiles, win)
            
            // Draw sprites
            shipSprite.Draw(win, pixel.IM.Moved(ship.pos))
        }
        // Update window
		win.Update()
	}
}


// [Update the ship velocity and position]
func updateShip(ship physObj, inputDirection pixel.Vec) (physObj) {
    // Give velocity a minimum limit or apply friction to the velocity if there is any
    if ship.vel.Len() <= 0.01 {
        ship.vel = pixel.ZV
    } else {
        ship.vel = ship.vel.Scaled(ship.frc)
    }
    
    // Add new velocity if there is input
    if inputDirection != pixel.ZV {
        ship.vel = ship.vel.Add(inputDirection.Scaled(ship.acc/inputDirection.Len()))
    }
    
    // Enforce soft boundary on ship
    boundsCollisions := inBounds(ship.pos, WINSIZE, BOUNDARY_RANGE)
    if boundsCollisions != pixel.ZV {
        if boundsCollisions.Y == 1 {
            ship.vel.Y -= borderForce(ship.acc, BOUNDARY_RANGE[0], WINSIZE.Y - ship.pos.Y)
        } else if boundsCollisions.Y == -1 {
            ship.vel.Y += borderForce(ship.acc, BOUNDARY_RANGE[1], ship.pos.Y)
        }
        if boundsCollisions.X == 1 {
            ship.vel.X -= borderForce(ship.acc, BOUNDARY_RANGE[2], WINSIZE.X - ship.pos.X)
        } else if boundsCollisions.X == -1 {
            ship.vel.X += borderForce(ship.acc, BOUNDARY_RANGE[3], ship.pos.X)
        }
    }
    
    // Add new velocity to the position if there is any input
    if ship.vel.Len() != 0 {ship.pos = ship.pos.Add(ship.vel)}
    
    return ship
}

// [Border force formula]
func borderForce(acc float64, boundaryRange float64, pos float64) (float64) {
    return acc * BOUNDARY_STRENGTH * (1 - pos / boundaryRange)
}

// [Check if pos is in bounds]
func inBounds(pos pixel.Vec, max pixel.Vec, bounds [4]float64) (pixel.Vec) {
    var boundCollision pixel.Vec = pixel.ZV
    if pos.Y >= max.Y - bounds[0] {
        boundCollision.Y = 1
    } else if pos.Y <= bounds[1] {
        boundCollision.Y = -1
    }
    if pos.X >= max.X - bounds[2] {
        boundCollision.X = 1
    } else if pos.X <= bounds[3] {
        boundCollision.X = -1
    }
    
    return boundCollision
}

// [Create a bullet if a slot is free]
func createBullet(bullets *[BULLET_ALLOC_SIZE]projectile, shipPos pixel.Vec) {
    // Loop through bullet array
    for i:=0;i<BULLET_ALLOC_SIZE;i++ {
        if bullets[i].loaded == false {
            bullets[i] = projectileTypes[0]
            bullets[i].phys.pos = bullets[i].phys.pos.Add(shipPos)
            
            _, count := bulletsWithinRadius(bullets, bullets[i].phys.pos, 30)
            println(count)
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

// [Handle user input for a single frame]
func handleInput(win *pixelgl.Window) (pixel.Vec, bool, bool) {
    var dirVec pixel.Vec = pixel.ZV
    var shootButton bool = false
    var pauseButton bool = false
         
    if win.Pressed(pixelgl.KeyUp)    {dirVec = dirVec.Add(inputVec[0])}
    if win.Pressed(pixelgl.KeyLeft)  {dirVec = dirVec.Add(inputVec[1])}
    if win.Pressed(pixelgl.KeyDown)  {dirVec = dirVec.Add(inputVec[2])}
    if win.Pressed(pixelgl.KeyRight) {dirVec = dirVec.Add(inputVec[3])}
    if win.Pressed(pixelgl.KeySpace) {shootButton = true}
    if win.Pressed(pixelgl.KeyEscape) {pauseButton = true}
    
    return dirVec, shootButton, pauseButton
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

// Lonely Main Function :(
func main() {pixelgl.Run(run)}