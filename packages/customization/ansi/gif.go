package ansi

import (
	"bytes"
	"fmt"
	"image/gif"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func (a *Ansi) GIF(path string, width, height int, delay int) error {
	dirPath, err := gifToImages(path)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	fileCount := 0
	for _, file := range files {
		if !file.IsDir() {
			fileCount++
		}
	}

	for i := 0; i < fileCount; i++ {
		buffer := new(bytes.Buffer)
		buffer.WriteString("\x1b[;H")
		_, err := NewAnsi(buffer).Image(fmt.Sprintf("%s/frame%d.png", dirPath, i), width, height)
		if err != nil {
			return err
		}

		_, err = a.Writer.Write(buffer.Bytes())
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	os.RemoveAll(dirPath)

	return nil
}

func gifToImages(path string) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	uniqueDirName := fmt.Sprintf("gif%d", r.Intn(100000))
	dir, err := os.MkdirTemp("", uniqueDirName)
	if err != nil {
		return dir, err
	}

	gifFile, err := os.Open(path)
	if err != nil {
		return dir, err
	}

	defer gifFile.Close()

	gifImage, err := gif.DecodeAll(gifFile)
	if err != nil {
		return dir, err
	}

	for i, frame := range gifImage.Image {
		pngFileName := fmt.Sprintf("%s/frame%d.png", dir, i)
		pngFile, err := os.Create(pngFileName)
		if err != nil {
			return dir, err
		}

		defer pngFile.Close()

		err = png.Encode(pngFile, frame)
		if err != nil {
			return dir, err
		}
	}

	return dir, nil
}
