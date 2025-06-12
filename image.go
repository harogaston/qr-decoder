package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type ImageRequest struct {
	Scale  int
	Pixels [][]color.Color
}

func WriteImage(req ImageRequest) {
	r := image.Rect(0, 0, len(req.Pixels)*req.Scale, len(req.Pixels[0])*req.Scale)
	img := image.NewRGBA(r)

	for y, row := range req.Pixels {
		for x, c := range row {
			writePix(img, y, x, req.Scale, c)
		}
	}

	file, err := os.Create("qr.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func writePix(img *image.RGBA, row, col, size int, color color.Color) {
	for i := range size {
		for j := range size {
			img.Set(col*size+i, row*size+j, color)
		}
	}
}
