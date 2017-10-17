package main

import (
	"github.com/choo8/opengl-tutorials-go/common"
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
	window, err := glfw.CreateWindow(1024, 768, "Tutorial 02 - Red Triangle", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	// Ensure we can capture the escape key being pressed below
	window.SetInputMode(glfw.StickyKeysMode, glfw.True)

	// Dark blue background
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	var vertexArrayId uint32
	gl.GenVertexArrays(1, &vertexArrayId)
	gl.BindVertexArray(vertexArrayId)

	// Create and compile our GLSL program from the shaders
	programId, err := common.LoadShaders("SimpleVertexShader.vertexshader", "SimpleFragmentShader.fragmentshader")

	vertexBufferData := []float32{
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		0.0, 1.0, 0.0,
	}

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBufferData)*4, gl.Ptr(vertexBufferData), gl.STATIC_DRAW)

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {
		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Use our shader
		gl.UseProgram(programId)

		// 1st attribute buffer : vertices
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

		// Draw the triangles!
		gl.DrawArrays(gl.TRIANGLES, 0, 3) // 3 indices starting 0 -> 1 triangle
		gl.DisableVertexAttribArray(0)

		// Swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// Cleanup VBO
	gl.DeleteBuffers(1, &vertexBuffer)
	gl.DeleteVertexArrays(1, &vertexArrayId)
	gl.DeleteProgram(programId)

	return
}
