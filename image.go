package main

import (
	"image/color"

	"github.com/fogleman/gg"
)

type Shape string

const (
	ShapeSquare  Shape = "square"
	ShapeCircle  Shape = "circle"
	ShapeRounded Shape = "rounded"
	ShapeDiamond Shape = "diamond"
)

type ImageRequest struct {
	Scale        int
	Pixels       [][]color.Color
	Shape        Shape
	BorderRadius float64 // For rounded rect (0-1, relative to module size)
}

func WriteImage(req ImageRequest) {
	rows := len(req.Pixels)
	cols := len(req.Pixels[0])
	width := cols * req.Scale
	height := rows * req.Scale

	dc := gg.NewContext(width, height)

	// Background (usually white, but let's assume transparent or white based on pixels)
	// We can fill with white first if needed.
	dc.SetColor(color.White)
	dc.Clear()

	for y, row := range req.Pixels {
		for x, c := range row {
			// Skip white/transparent modules if we cleared with white?
			// Or just draw everything.
			// Let's draw only non-white modules to allow background to show?
			// But pixels contain color.

			// If color is white, we might want to skip drawing if we cleared with white,
			// but for shapes like circles, the background matters.
			// If we draw a white circle on white background, it's invisible.
			// If we draw a black circle, it's visible.

			// Let's assume we draw every module with its color.
			dc.SetColor(c)

			posX := float64(x * req.Scale)
			posY := float64(y * req.Scale)
			size := float64(req.Scale)

			// Padding for shapes to not touch?
			// Usually QR modules touch.
			// For circles/rounded, they might not fill the square fully.

			switch req.Shape {
			case ShapeCircle:
				radius := size / 2
				dc.DrawCircle(posX+radius, posY+radius, radius)
				dc.Fill()
			case ShapeRounded:
				// Add 10% total padding (5% on each side)
				padding := size * 0.05
				newSize := size - 2*padding

				// BorderRadius is relative to size? Or absolute?
				// Let's say req.BorderRadius is 0.0 to 0.5 (radius/size).
				// If not provided, default to 0.2?
				r := req.BorderRadius * newSize
				if r == 0 {
					r = newSize * 0.4 // Default
				}
				dc.DrawRoundedRectangle(posX+padding, posY+padding, newSize, newSize, r)
				dc.Fill()
			case ShapeDiamond:
				// Rotate 45 degrees around center
				dc.Push()
				dc.RotateAbout(gg.Radians(45), posX+size/2, posY+size/2)
				// Draw square (which becomes diamond)
				// We need to scale it down to fit?
				// A rotated square of size S has bounding box S*sqrt(2).
				// To fit in S, we need to scale by 1/sqrt(2) ~ 0.707.
				scale := 1.0 / 1.41421356
				offset := (size - size*scale) / 2
				dc.DrawRectangle(posX+offset, posY+offset, size*scale, size*scale)
				dc.Fill()
				dc.Pop()
			default: // ShapeSquare
				dc.DrawRectangle(posX, posY, size, size)
				dc.Fill()
			}
		}
	}

	dc.SavePNG("qr.png")
}
