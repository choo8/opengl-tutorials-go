package common

import (
	"encoding/binary"
	"fmt"
	"github.com/go-gl/gl/v4.5-core/gl"
	"log"
	"os"
	//"strings"
)

const (
	FOURCC_DXT1 string = "0x31545844"
	FOURCC_DXT3 string = "0x33545844"
	FOURCC_DXT5 string = "0x35545844"
)

func LoadBMPCustom(imagepath string) uint32 {
	// Data read from the header of the BMP file
	header := make([]byte, 54)
	var dataPos uint32
	var imageSize uint32
	var width, height uint32
	// Actual RGB data
	var data []byte

	// Open the file
	f, err := os.Open(imagepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read the header, i.e. the 54 first bytes

	// If less than 54 bytes are read, problem
	f.Read(header)
	// A BMP files always begins with "BM"
	if string(header[0]) != "B" || string(header[1]) != "M" {
		fmt.Println("Not a correct BMP file")
		return 0
	}
	// Make sure this is a 24bpp file
	if binary.LittleEndian.Uint32(header[30:34]) != 0 {
		fmt.Println("Not a correct BMP file")
	}
	if binary.LittleEndian.Uint32(header[28:32]) != 24 {
		fmt.Println("Not a correct BMP file")
	}

	// Read the information about the image
	dataPos = binary.LittleEndian.Uint32(header[10:14])
	imageSize = binary.LittleEndian.Uint32(header[34:38])
	width = binary.LittleEndian.Uint32(header[18:22])
	height = binary.LittleEndian.Uint32(header[22:26])

	// Some BMP files are misformatted, guess missing information
	if imageSize == 0 {
		imageSize = width * height * 3 // 3 : one byte for each Red, Green and Blue component
	}
	if dataPos == 0 {
		dataPos = 54 // The BMP header is done that way
	}

	// Create a buffer
	data = make([]byte, imageSize)

	// Read the actual data from the file into the buffer
	f.Read(data)

	// Create one OpenGL texture
	var textureId uint32
	gl.GenTextures(1, &textureId)

	// "Bind" the newly created texture : all future texture functions will modify this texture
	gl.BindTexture(gl.TEXTURE_2D, textureId)

	// Give the image to OpenGL
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.BGR, gl.UNSIGNED_BYTE, gl.Ptr(data))

	// ... nice trilinear filtering ...
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	// ... which requires mipmaps. Generate them automatically.
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return textureId
}

/*
func LoadDDS(imagepath string) uint32 {
	// try to open the file
	f, err := os.Open(imagepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// verify the type of file
	filecode := make([]byte, 4)
	f.Read(filecode)
	if strings.Compare(string(filecode), "DDS ") != 0 {
		return 0
	}

	// get surface desc
	header := make([]byte, 124)
	f.Read(header)

	height := binary.LittleEndian.Uint32(header[8:12])
	width := binary.LittleEndian.Uint32(header[12:16])
	linearSize := binary.LittleEndian.Uint32(header[16:20])
	mipMapCount := binary.LittleEndian.Uint32(header[24:28])
	fourCC := binary.LittleEndian.Uint32(header[80:84])

	var bufsize uint32
	buffer := make([]byte, bufsize)
	// how big is it going to be including all mipmaps?
	if mipMapCount > 1 {
		bufsize := linearSize * 2
	} else {
		bufsize := linearSize
	}
	f.Read(buffer)

	return 0
}
*/
