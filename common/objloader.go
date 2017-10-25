package common

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"io"
	"log"
	"os"
	"strings"
)

func LoadOBJ(path string, outVertices []mgl32.Vec3, outUvs []mgl32.Vec2, outNormals []mgl32.Vec3) ([]mgl32.Vec3, []mgl32.Vec2, []mgl32.Vec3, bool) {
	var vertexIndices, uvIndices, normalIndices []uint32
	var tempVertices, tempNormals []mgl32.Vec3
	var tempUvs []mgl32.Vec2

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for {
		var lineHeader string
		// read the first word of the line
		_, err := fmt.Fscanf(f, "%s", &lineHeader)
		if err == io.EOF {
			break // EOF = End Of File. Quit the loop.
		}

		// else : parse lineHeader

		if strings.Compare(string(lineHeader), "v") == 0 {
			var vertex mgl32.Vec3
			fmt.Fscanf(f, "%f %f %f\n", &vertex[0], &vertex[1], &vertex[2])
			tempVertices = append(tempVertices, vertex)
		} else if strings.Compare(string(lineHeader), "vt") == 0 {
			var uv mgl32.Vec2
			fmt.Fscanf(f, "%f %f\n", &uv[0], &uv[1])
			uv[1] = -uv[1] // Invert V coordinate since we will only use DDS texture, which are inverted. Remove if you want to use TGA or BMP loaders.
			tempUvs = append(tempUvs, uv)
		} else if strings.Compare(string(lineHeader), "vn") == 0 {
			var normal mgl32.Vec3
			fmt.Fscanf(f, "%f %f %f\n", &normal[0], &normal[1], &normal[2])
			tempNormals = append(tempNormals, normal)
		} else if strings.Compare(string(lineHeader), "f") == 0 {
			//var vertex1, vertex2, vertex3 string
			var vertexIndex, uvIndex, normalIndex [3]uint32
			matches, _ := fmt.Fscanf(f, "%d/%d/%d %d/%d/%d %d/%d/%d\n", &vertexIndex[0], &uvIndex[0], &normalIndex[0], &vertexIndex[1], &uvIndex[1], &normalIndex[1], &vertexIndex[2], &uvIndex[2], &normalIndex[2])
			if matches != 9 {
				fmt.Println("File can't be read by our simple parser :-( Try exporting with other options")
				return []mgl32.Vec3{}, []mgl32.Vec2{}, []mgl32.Vec3{}, false
			}
			vertexIndices = append(vertexIndices, vertexIndex[0])
			vertexIndices = append(vertexIndices, vertexIndex[1])
			vertexIndices = append(vertexIndices, vertexIndex[2])
			uvIndices = append(uvIndices, uvIndex[0])
			uvIndices = append(uvIndices, uvIndex[1])
			uvIndices = append(uvIndices, uvIndex[2])
			normalIndices = append(normalIndices, normalIndex[0])
			normalIndices = append(normalIndices, normalIndex[1])
			normalIndices = append(normalIndices, normalIndex[2])
		}
	}

	// For each vertex of each triangle
	for i := 0; i < len(vertexIndices); i++ {

		// Get the indices of its attributes
		vertexIndex := vertexIndices[i]
		uvIndex := uvIndices[i]
		normalIndex := normalIndices[i]

		// Get the attributes thanks to the index
		vertex := tempVertices[vertexIndex-1]
		uv := tempUvs[uvIndex-1]
		normal := tempNormals[normalIndex-1]

		// Put the attributes in buffers
		outVertices = append(outVertices, vertex)
		outUvs = append(outUvs, uv)
		outNormals = append(outNormals, normal)
	}

	return outVertices, outUvs, outNormals, true
}
