package common

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"sync"
)

var once sync.Once
var lastTime float64

var ViewMatrix, ProjectionMatrix mgl32.Mat4
var position mgl32.Vec3
var horizontalAngle, verticalAngle, initialFoV, speed, mouseSpeed float64

func init() {
	// Initial position : on +Z
	position = mgl32.Vec3{0, 0, 5}
	// Initial horizontal angle : toward -Z
	horizontalAngle = 3.14
	// Initial vertical angle : none
	verticalAngle = 0.0
	// Initial Field of View
	initialFoV = 45.0

	speed = 3.0 // 3 units / second
	mouseSpeed = 0.005
}

func ComputeMatricesFromInputs(window *glfw.Window) {
	// glfwGetTime is called only once, the first time this function is called
	getLastTime := func() {
		lastTime = glfw.GetTime()
	}
	once.Do(getLastTime)

	// Compute time difference between current and last frame
	currentTime := glfw.GetTime()
	deltaTime := float32(currentTime - lastTime)

	// Get mouse position
	var xPos, yPos float64
	xPos, yPos = window.GetCursorPos()

	// Reset mouse position for next frame
	window.SetCursorPos(1024/2, 768/2)

	// Compute new orientation
	horizontalAngle += mouseSpeed * float64(1024/2-xPos)
	verticalAngle += mouseSpeed * float64(768/2-yPos)

	// Direction : Spherical coordinates to Cartesian coordinates conversion
	direction := mgl32.Vec3{float32(math.Cos(verticalAngle) * math.Sin(horizontalAngle)), float32(math.Sin(verticalAngle)), float32(math.Cos(verticalAngle) * math.Cos(horizontalAngle))}

	// Right vector
	right := mgl32.Vec3{float32(math.Sin(horizontalAngle - 3.14/2.0)), 0, float32(math.Cos(horizontalAngle - 3.14/2.0))}

	// Up vector
	up := right.Cross(direction)

	// Move forward
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		position = position.Add(direction.Mul(deltaTime * float32(speed)))
	}

	// Move backward
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		position = position.Sub(direction.Mul(deltaTime * float32(speed)))
	}

	// Strafe right
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		position = position.Add(right.Mul(deltaTime * float32(speed)))
	}

	// Strafe left
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		position = position.Sub(right.Mul(deltaTime * float32(speed)))
	}

	FoV := initialFoV // - 5 * glfwGetMouseWheel(); // Now GLFW 3 requires setting up a callback for this. It's a bit too complicated for this beginner's tutorial, so it's disabled instead.

	// Projection matrix : 45° Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	ProjectionMatrix = mgl32.Perspective(mgl32.DegToRad(float32(FoV)), float32(1024)/float32(768), 0.1, 100.0)
	// Camera matrix
	ViewMatrix = mgl32.LookAtV(position, position.Add(direction), up)

	// For the next frame, the "last time" will be "now"
	lastTime = currentTime
}
