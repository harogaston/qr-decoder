package writers

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	svg "github.com/ajstarks/svgo"
)

const path string = "qr.svg"
const roundedCornerPath = "M${x},${y} v ${size} h ${size} v ${-size/2} a ${size/2},${size/2} 0,0,0 ${-size/2},${-size/2} z"

type SVGRequest struct {
	Scale  int
	Pixels [][]color.Color
	Shape  Shape
}

func WriteSVG(req SVGRequest) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating SVG file:", err)
	}
	defer file.Close()

	dim := len(req.Pixels)
	size := dim * req.Scale
	canvas := svg.New(file)
	canvas.Start(size, size)

	for y, row := range req.Pixels {
		for x, c := range row {
			switch req.Shape {
			case ShapeCircle:
				canvas.Circle(x*req.Scale+req.Scale/2, y*req.Scale+req.Scale/2, req.Scale/2, CircleStyle(c, 0))
				// Superimpose finder patterns
			case ShapeSquare:
				canvas.Square(x*req.Scale, y*req.Scale, req.Scale, SquareStyle(c))
			case ShapeSlanted:
				pathExpression := strings.ReplaceAll(roundedCornerPath, "${x}", fmt.Sprintf("%d", x*req.Scale))
				pathExpression = strings.ReplaceAll(pathExpression, "${y}", fmt.Sprintf("%d", y*req.Scale))
				pathExpression = strings.ReplaceAll(pathExpression, "${size}", fmt.Sprintf("%d", req.Scale))
				pathExpression = strings.ReplaceAll(pathExpression, "${-size/2}", fmt.Sprintf("%d", -req.Scale/2))
				pathExpression = strings.ReplaceAll(pathExpression, "${size/2}", fmt.Sprintf("%d", req.Scale/2))
				canvas.Path(pathExpression, SquareStyle(c))
			default:
				fmt.Println("Not implemented shape:", req.Shape)
			}
		}
	}
	// Superimpose finder patterns
	switch req.Shape {
	case ShapeCircle:
		CircleFinderPattern(canvas, 4, 4, req.Scale)
		CircleFinderPattern(canvas, (dim - 7 - 4), 4, req.Scale)
		CircleFinderPattern(canvas, 4, (dim - 7 - 4), req.Scale)
	case ShapeSquare:
	case ShapeSlanted:
	}
	canvas.End()
}

func CircleFinderPattern(canvas *svg.SVG, x, y int, scale int) {
	canvas.Square(x*scale, y*scale, 8*scale, SquareStyle(color.White))
	canvas.Circle(
		int((float32(x)+float32(3.5))*float32(scale)),
		int((float32(y)+float32(3.5))*float32(scale)),
		int(float32(3.5)*float32(scale)), CircleStyle(color.Black, 4),
	)
	canvas.Circle(
		int((float32(x)+float32(3.5))*float32(scale)),
		int((float32(y)+float32(3.5))*float32(scale)),
		int(float32(2.5)*float32(scale)), CircleStyle(color.White, 4),
	)
	canvas.Circle(
		int((float32(x)+float32(3.5))*float32(scale)),
		int((float32(y)+float32(3.5))*float32(scale)),
		int(float32(1.5)*float32(scale)), CircleStyle(color.Black, 4),
	)
	// canvas.Circle((x+3.5)*scale, (y+3.5)*scale, int(float32(2.5)*float32(scale)), CircleStyle(color.White, 4))
	// canvas.Circle((x+3.5)*scale, (y+3.5)*scale, int(float32(1.5)*float32(scale)), CircleStyle(color.Black, 4))
}

func SquareStyle(c color.Color) string {
	if c == color.Black {
		return "fill:black;stroke:none"
	}
	return "fill:white;stroke:none"
}

func CircleStyle(c color.Color, border int) string {
	if border <= 0 {
		border = 2
	}
	if c == color.Black {
		return fmt.Sprintf("fill:black;stroke:white;stroke-width:%d", border)
	}
	return fmt.Sprintf("fill:white;stroke:white;stroke-width:%d", border)
}
