package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Handle user input for a single frame
func handleInput(win *pixelgl.Window) {
	input.dir = pixel.Vec{ // Initialize dirVec with joystick pos
		X: win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftX),
		Y: win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftY) * -1,
	}

	// Add keyboard input to the direction vector
	if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
		input.dir = input.dir.Add(inputLookup[0])
	}
	if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
		input.dir = input.dir.Add(inputLookup[1])
	}
	if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
		input.dir = input.dir.Add(inputLookup[2])
	}
	if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
		input.dir = input.dir.Add(inputLookup[3])
	}

	// Trim directional vector and add deadzone
	if dirVecLen := input.dir.Len(); dirVecLen > 1 {
		input.dir = input.dir.Scaled(1 / dirVecLen)
	} else if dirVecLen < AXIS_DEADZONE {
		input.dir = pixel.ZV
	}

	// Shoot
	if win.JoystickPressed(pixelgl.Joystick1, pixelgl.ButtonRightBumper) || win.Pressed(pixelgl.MouseButtonLeft) {
		input.shoot = true
	} else {
		input.shoot = false
	}

	// Roll
	if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonA) || win.JustPressed(pixelgl.MouseButtonRight) {
		input.roll = true
	} else {
		input.roll = false
	}
}

// Button functions

func pauseButton() bool {
	return win.JustPressed(pixelgl.KeyEscape) || win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonStart)
}

func startButton() bool {
	return win.JustPressed(pixelgl.KeySpace) || win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonStart)
}
