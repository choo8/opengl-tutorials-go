package main

import (
	"github.com/choo8/opengl-tutorials-go/common"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
	window, err := glfw.CreateWindow(1024, 768, "Tutorial 05 - Textured Cube", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	// Ensure we can capture the escape key being pressed below
	window.SetInputMode(glfw.StickyKeysMode, glfw.True)

	// Dark blue background
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	// Enable depth test
	gl.Enable(gl.DEPTH_TEST)
	// Accept fragment if it closer to the camera than the former one
	gl.DepthFunc(gl.LESS)

	var vertexArrayId uint32
	gl.GenVertexArrays(1, &vertexArrayId)
	gl.BindVertexArray(vertexArrayId)

	// Create and compile our GLSL program from the shaders
	programId, err := common.LoadShaders("TransformVertexShader.vertexshader", "TextureFragmentShader.fragmentshader")

	// Get a handle for our "MVP" uniform
	matrixId := gl.GetUniformLocation(programId, gl.Str("MVP"+"\x00"))

	// Projection matrix : 45 degrees Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(1024)/float32(768), 0.1, 100.0)
	// Camera matrix
	view := mgl32.LookAtV(mgl32.Vec3{4, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	// Model matrix : an identity matrix (model will be at the origin)
	model := mgl32.Ident4()
	// Our ModelViewProjection : multiplication of our 3 matrices
	MVP := projection.Mul4(view.Mul4(model))

	// Load the texture using any two methods
	//texture := common.LoadBMPCustom("uvtemplate.bmp")
	texture := common.LoadDDS("uvtemplate.DDS")

	// Get a handle for our "myTextureSampler" uniform
	textureId := gl.GetUniformLocation(programId, gl.Str("myTextureSampler"+"\x00"))

	// Our vertices. Tree consecutive floats give a 3D vertex; Three consecutive vertices give a triangle.
	// A cube has 6 faces with 2 triangles each, so this makes 6*2=12 triangles, and 12*3 vertices
	gVertexBufferData := []float32{
		-1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		0.0, 1.0, 0.0,
	}

	// Two UV coordinatesfor each vertex. They were created with Blender.
	gUvBufferData := []float32{
		0.000059, 1.0 - 0.000004,
		0.000103, 1.0 - 0.336048,
		0.335973, 1.0 - 0.335903,
		1.000023, 1.0 - 0.000013,
		0.667979, 1.0 - 0.335851,
		0.999958, 1.0 - 0.336064,
		0.667979, 1.0 - 0.335851,
		0.336024, 1.0 - 0.671877,
		0.667969, 1.0 - 0.671889,
		1.000023, 1.0 - 0.000013,
		0.668104, 1.0 - 0.000013,
		0.667979, 1.0 - 0.335851,
		0.000059, 1.0 - 0.000004,
		0.335973, 1.0 - 0.335903,
		0.336098, 1.0 - 0.000071,
		0.667979, 1.0 - 0.335851,
		0.335973, 1.0 - 0.335903,
		0.336024, 1.0 - 0.671877,
		1.000004, 1.0 - 0.671847,
		0.999958, 1.0 - 0.336064,
		0.667979, 1.0 - 0.335851,
		0.668104, 1.0 - 0.000013,
		0.335973, 1.0 - 0.335903,
		0.667979, 1.0 - 0.335851,
		0.335973, 1.0 - 0.335903,
		0.668104, 1.0 - 0.000013,
		0.336098, 1.0 - 0.000071,
		0.000103, 1.0 - 0.336048,
		0.000004, 1.0 - 0.671870,
		0.336024, 1.0 - 0.671877,
		0.000103, 1.0 - 0.336048,
		0.336024, 1.0 - 0.671877,
		0.335973, 1.0 - 0.335903,
		0.667969, 1.0 - 0.671889,
		1.000004, 1.0 - 0.671847,
		0.667979, 1.0 - 0.335851,
	}

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(gVertexBufferData)*4, gl.Ptr(gVertexBufferData), gl.STATIC_DRAW)

	var uvBuffer uint32
	gl.GenBuffers(1, &uvBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, uvBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(gUvBufferData)*4, gl.Ptr(gUvBufferData), gl.STATIC_DRAW)

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {
		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Use our shader
		gl.UseProgram(programId)

		// Send our transformation to the currently bound shader,
		// in the "MVP" uniform
		gl.UniformMatrix4fv(matrixId, 1, false, &MVP[0])

		// Bind our texture in Texture Unit 0
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		// Set our "myTextureSampler" sampler to use Texture Unit 0
		gl.Uniform1i(textureId, 0)

		// 1st attribute buffer : vertices
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

		// 2nd attribute buffer : UVs
		gl.EnableVertexAttribArray(1)
		gl.BindBuffer(gl.ARRAY_BUFFER, uvBuffer)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)

		// Draw the triangles!
		gl.DrawArrays(gl.TRIANGLES, 0, 12*3) // 3 indices starting 0 -> 1 triangle
		gl.DisableVertexAttribArray(0)
		gl.DisableVertexAttribArray(1)

		// Swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// Cleanup VBO
	gl.DeleteBuffers(1, &vertexBuffer)
	gl.DeleteBuffers(1, &uvBuffer)
	gl.DeleteProgram(programId)
	gl.DeleteTextures(1, &texture)
	gl.DeleteVertexArrays(1, &vertexArrayId)

	return
}
