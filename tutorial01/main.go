package main

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
)

func init() {
	runtime.LockOSThread()

	// Initialize GL
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Initialize GLFW
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// Open a window and create its OpenGL context
	window, err := glfw.CreateWindow(1024, 768, "Tutorial 01", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	// Ensure we can capture the escape key being pressed below
	window.SetInputMode(glfw.StickyKeysMode, glfw.True)

	// Dark blue background
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {
		// Clear the screen. It's not mentioned before Tutorial 02, but it can cause flickering, so it's there nonetheless.
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw nothing, see you in tutorial 2!

		// Swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}

	return
}
