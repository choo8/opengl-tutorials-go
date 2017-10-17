package common

import (
	"bufio"
	"fmt"
	"github.com/go-gl/gl/v4.5-core/gl"
	"log"
	"os"
	"strings"
)

func LoadShaders(vertexFilePath, fragmentFilePath string) (uint32, error) {
	// Create the shaders
	vertexShaderId := gl.CreateShader(gl.VERTEX_SHADER)
	fragmentShaderId := gl.CreateShader(gl.FRAGMENT_SHADER)

	// Read the Vertex Shader code from the file
	vertexShaderCode := ""
	vertexF, err := os.Open(vertexFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer vertexF.Close()

	vertexShaderStream := bufio.NewScanner(vertexF)
	for vertexShaderStream.Scan() {
		vertexShaderCode += "\n" + vertexShaderStream.Text()
	}
	if err := vertexShaderStream.Err(); err != nil {
		log.Fatal(err)
	}
	vertexShaderCode += "\x00"

	// Read the Fragment Shader code from the file
	fragmentShaderCode := ""
	fragmentF, err := os.Open(fragmentFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fragmentF.Close()

	fragmentShaderStream := bufio.NewScanner(fragmentF)
	for fragmentShaderStream.Scan() {
		fragmentShaderCode += "\n" + fragmentShaderStream.Text()
	}
	if err := fragmentShaderStream.Err(); err != nil {
		log.Fatal(err)
	}
	fragmentShaderCode += "\x00"

	var result int32
	var infoLogLength int32

	// Compile Vertex Shader
	vertexSourcePointer, free := gl.Strs(vertexShaderCode)
	defer free()
	gl.ShaderSource(vertexShaderId, 1, vertexSourcePointer, nil)
	gl.CompileShader(vertexShaderId)

	// Check Vertex Shader
	gl.GetShaderiv(vertexShaderId, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(vertexShaderId, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		vertexShaderErrorMessage := strings.Repeat("\x00", int(infoLogLength+1))
		gl.GetShaderInfoLog(vertexShaderId, infoLogLength, nil, gl.Str(vertexShaderErrorMessage))
		fmt.Println(vertexShaderErrorMessage)
	}

	// Compile Fragment Shader
	fragmentSourcePointer, free := gl.Strs(fragmentShaderCode)
	defer free()
	gl.ShaderSource(fragmentShaderId, 1, fragmentSourcePointer, nil)
	gl.CompileShader(fragmentShaderId)

	// Check Fragment Shader
	gl.GetShaderiv(fragmentShaderId, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(fragmentShaderId, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		fragmentShaderErrorMessage := strings.Repeat("\x00", int(infoLogLength+1))
		gl.GetShaderInfoLog(fragmentShaderId, infoLogLength, nil, gl.Str(fragmentShaderErrorMessage))
		fmt.Println(fragmentShaderErrorMessage)
	}

	// Link the program
	programId := gl.CreateProgram()
	gl.AttachShader(programId, vertexShaderId)
	gl.AttachShader(programId, fragmentShaderId)
	gl.LinkProgram(programId)

	// Check the program
	gl.GetProgramiv(programId, gl.LINK_STATUS, &result)
	gl.GetProgramiv(programId, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		programErrorMessage := strings.Repeat("\x00", int(infoLogLength+1))
		gl.GetProgramInfoLog(programId, infoLogLength, nil, gl.Str(programErrorMessage))
		fmt.Println(programErrorMessage)
	}

	gl.DetachShader(programId, vertexShaderId)
	gl.DetachShader(programId, fragmentShaderId)

	gl.DeleteShader(vertexShaderId)
	gl.DeleteShader(fragmentShaderId)

	return programId, nil
}
