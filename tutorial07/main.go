package main

import (
	"github.com/choo8/opengl-tutorials-go/common"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"runtime"
	"unsafe"
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
	window, err := glfw.CreateWindow(1024, 768, "Tutorial 07 - Model Loading", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	// Ensure we can capture the escape key being pressed below
	window.SetInputMode(glfw.StickyKeysMode, glfw.True)
	// Hide the mouse and enable unlimited mouvement
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Set the mouse at the center of the screen
	glfw.PollEvents()
	window.SetCursorPos(1024/2, 768/2)

	// Dark blue background
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	// Enable depth test
	gl.Enable(gl.DEPTH_TEST)
	// Accept fragment if it closer to the camera than the former one
	gl.DepthFunc(gl.LESS)

	// Cull triangles which normal is not towards the camera
	gl.Enable(gl.CULL_FACE)

	var vertexArrayId uint32
	gl.GenVertexArrays(1, &vertexArrayId)
	gl.BindVertexArray(vertexArrayId)

	// Create and compile our GLSL program from the shaders
	programId, err := common.LoadShaders("TransformVertexShader.vertexshader", "TextureFragmentShader.fragmentshader")

	// Get a handle for our "MVP" uniform
	matrixId := gl.GetUniformLocation(programId, gl.Str("MVP"+"\x00"))

	// Load the texture using any two methods
	texture := common.LoadDDS("uvmap.DDS")

	// Get a handle for our "myTextureSampler" uniform
	textureId := gl.GetUniformLocation(programId, gl.Str("myTextureSampler"+"\x00"))

	// Read our .obj file
	var vertices, normals []mgl32.Vec3
	var uvs []mgl32.Vec2
	vertices, uvs, normals, _ = common.LoadOBJ("cube.obj", vertices, uvs, normals)

	// Load it into a VBO
	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(mgl32.Vec3{})), gl.Ptr(vertices), gl.STATIC_DRAW)

	var uvBuffer uint32
	gl.GenBuffers(1, &uvBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, uvBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(uvs)*int(unsafe.Sizeof(mgl32.Vec2{})), gl.Ptr(uvs), gl.STATIC_DRAW)

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {
		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Use our shader
		gl.UseProgram(programId)

		// Compute the MVP matrix from keyboard and mouse input
		common.ComputeMatricesFromInputs(window)
		model := mgl32.Ident4()
		MVP := common.ProjectionMatrix.Mul4(common.ViewMatrix.Mul4(model))

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

	// Cleanup VBO and shader
	gl.DeleteBuffers(1, &vertexBuffer)
	gl.DeleteBuffers(1, &uvBuffer)
	gl.DeleteProgram(programId)
	gl.DeleteTextures(1, &texture)
	gl.DeleteVertexArrays(1, &vertexArrayId)

	return
}
