package ansi

import (
	"bytes"
	"fmt"
	"image"
	"os"
)

func (a *Ansi) Image(path string, width, height int) (int, error) {
	buffer := new(bytes.Buffer)

	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}

	buffer.WriteString("\x1b[0m")

	for y := 0; y < height-1; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x*img.Bounds().Dx()/width, y*img.Bounds().Dy()/height).RGBA()
			buffer.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm ", r>>8, g>>8, b>>8))
		}

		buffer.WriteString("\x1b[0m\r\n")
	}

	return a.Writer.Write(buffer.Bytes())
}
