package writers

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/harogaston/qr-decoder/images"
	// "github.com/harogaston/qr-decoder/images"
)

const path string = "qr.svg"
const roundedCornerPath = "M${x},${y} v ${size} h ${size} v ${-size/2} a ${size/2},${size/2} 0,0,0 ${-size/2},${-size/2} z"
const squircleCurvature = 0.75
const squirclePath = "M0,${size/2} C 0 ${curvA}, ${curvA} 0, ${size/2} 0 S ${size} ${curvA}, ${size} ${size/2}, ${curvB} ${size}, ${size/2} ${size}, 0 ${curvB}, 0 ${size/2}"

// M 0 80
// C 0 20, 20 0, 80 0
// S 160 20, 160 80, 140 160
// 80 160, 0 140, 0 80
const logoRelativeSize = 5

type SVGRequest struct {
	Scale  int
	Pixels [][]color.Color
	Shape  Shape
	Logo   string
	Color  color.Color
}

func WriteSVG(req SVGRequest) {
	if req.Color == nil {
		req.Color = color.RGBA{R: 0, G: 96, B: 250, A: 50} //rgb(0,96,250)
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating SVG file:", err)
	}
	defer file.Close()

	dim := len(req.Pixels)
	size := dim * req.Scale
	canvas := svg.New(file)
	canvas.Start(size, size)

	canvas.Def()
	// Cell
	switch req.Shape {
	case ShapeCircle:
		canvas.Circle(0, 0, req.Scale/2, `id="cell"`)
		// Finder pattern block
		canvas.Group(`id="finderpattern"`)
		CircleFinderPattern(canvas, req.Color, 4, 4, req.Scale)
		CircleFinderPattern(canvas, req.Color, (dim - 7 - 4), 4, req.Scale)
		CircleFinderPattern(canvas, req.Color, 4, (dim - 7 - 4), req.Scale)
		canvas.Gend()
	case ShapeSquare:
		canvas.Square(0, 0, req.Scale, `id="cell"`)
		canvas.Group(`id="finderpattern"`)
		SquareFinderPattern(canvas, req.Color, 4, 4, req.Scale)
		SquareFinderPattern(canvas, req.Color, (dim - 7 - 4), 4, req.Scale)
		SquareFinderPattern(canvas, req.Color, 4, (dim - 7 - 4), req.Scale)
		canvas.Gend()
	case ShapeSquircle:
		squirclePathExp := strings.ReplaceAll(squirclePath, "${size}", fmt.Sprintf("%d", req.Scale))
		squirclePathExp = strings.ReplaceAll(squirclePathExp, "${size/2}", fmt.Sprintf("%d", req.Scale/2))
		squirclePathExp = strings.ReplaceAll(squirclePathExp, "${curvA}", fmt.Sprintf("%.f", float32(req.Scale/2.)*float32(1.-squircleCurvature)))
		squirclePathExp = strings.ReplaceAll(squirclePathExp, "${curvB}", fmt.Sprintf("%.f", float32(req.Scale)-float32(req.Scale/2.)*float32(1.-squircleCurvature)))
		canvas.Path(squirclePathExp, `id="cell"`, `transform="rotate(0,80,80)"`)

		// Finder pattern block
		canvas.Group(`id="finderpattern"`)
		SquircleFinderPattern(canvas, req.Color, 4, 4, req.Scale)
		SquircleFinderPattern(canvas, req.Color, (dim - 7 - 4), 4, req.Scale)
		SquircleFinderPattern(canvas, req.Color, 4, (dim - 7 - 4), req.Scale)
		canvas.Gend()
	case ShapeSlanted:
		// pathExpression := strings.ReplaceAll(roundedCornerPath, "${x}", fmt.Sprintf("%d", x*req.Scale))
		// pathExpression = strings.ReplaceAll(pathExpression, "${y}", fmt.Sprintf("%d", y*req.Scale))
		// pathExpression = strings.ReplaceAll(pathExpression, "${size}", fmt.Sprintf("%d", req.Scale))
		// pathExpression = strings.ReplaceAll(pathExpression, "${-size/2}", fmt.Sprintf("%d", -req.Scale/2))
		// pathExpression = strings.ReplaceAll(pathExpression, "${size/2}", fmt.Sprintf("%d", req.Scale/2))
		// canvas.Path(pathExpression, SquareStyle(c))
	default:
		fmt.Println("Not implemented shape:", req.Shape)
	}

	canvas.DefEnd()

	for y, row := range req.Pixels {
		for x, c := range row {
			if c == color.Black {
				switch req.Shape {
				case ShapeSquircle:
					canvas.Use(x*req.Scale, y*req.Scale, "#cell", StrokeStyle(req.Color, c))
				case ShapeCircle:
					canvas.Use(x*req.Scale+req.Scale/2, y*req.Scale+req.Scale/2, "#cell", StrokeStyle(req.Color, c))
				case ShapeSquare:
					canvas.Square(x*req.Scale, y*req.Scale, req.Scale, NoStrokeStyle(req.Color, c))
				case ShapeSlanted:
					pathExpression := strings.ReplaceAll(roundedCornerPath, "${x}", fmt.Sprintf("%d", x*req.Scale))
					pathExpression = strings.ReplaceAll(pathExpression, "${y}", fmt.Sprintf("%d", y*req.Scale))
					pathExpression = strings.ReplaceAll(pathExpression, "${size}", fmt.Sprintf("%d", req.Scale))
					pathExpression = strings.ReplaceAll(pathExpression, "${-size/2}", fmt.Sprintf("%d", -req.Scale/2))
					pathExpression = strings.ReplaceAll(pathExpression, "${size/2}", fmt.Sprintf("%d", req.Scale/2))
					canvas.Path(pathExpression, NoStrokeStyle(req.Color, c))
				default:
					fmt.Println("Not implemented shape:", req.Shape)
				}
			}
		}
	}

	// Superimpose finder patterns
	switch req.Shape {
	case ShapeCircle:
		fallthrough
	case ShapeSquare:
		canvas.Use(0, 0, "#finderpattern")
	case ShapeSlanted:
	}

	// Superimpose logo if any
	if req.Logo != "" {
		// Create rounded logo version
		logoPath := images.MakeLogo(req.Logo)
		logoSize := dim / 5 * req.Scale

		// Delete overlapping QR modules
		// For all cells inside the square of logoSize, if any cell corners overlap, draw white square
		logoDim := dim/logoRelativeSize + 1 // logo size in cells

		startCell := (dim - logoDim) / 2
		endCell := startCell + logoDim

		for y := startCell; y <= endCell; y++ {
			for x := startCell; x <= endCell; x++ {
				// Check if inside logo circle area
				center := dim / 2
				radius := logoDim / 2
				radius += 1 // Add some padding
				dx := x - center
				dy := y - center
				distance := dx*dx + dy*dy
				if distance < radius*radius {
					canvas.Square(x*req.Scale, y*req.Scale, req.Scale, NoStrokeStyle(req.Color, color.White))
				}
			}
		}

		// Place logo at center
		logoPos := (size - logoSize) / 2
		canvas.Image(logoPos, logoPos, logoSize, logoSize, logoPath)
	}
	canvas.End()
}

func SquircleFinderPattern(canvas *svg.SVG, targetColor color.Color, x, y int, scale int) {
	canvas.Square(x*scale, y*scale, 8*scale, NoStrokeStyle(targetColor, color.White))
}

func SquareFinderPattern(canvas *svg.SVG, targetColor color.Color, x, y int, scale int) {
	canvas.Square(x*scale, y*scale, 8*scale, NoStrokeStyle(targetColor, color.White))
	canvas.Square((x)*scale, (y)*scale, 7*scale, StrokeStyle(targetColor, color.Black))
	canvas.Square((x+1)*scale, (y+1)*scale, 5*scale, StrokeStyle(targetColor, color.White))
	canvas.Square((x+2)*scale, (y+2)*scale, 3*scale, StrokeStyle(targetColor, color.Black))
}

func CircleFinderPattern(canvas *svg.SVG, targetColor color.Color, x, y int, scale int) {
	canvas.Square(x*scale, y*scale, 8*scale, NoStrokeStyle(targetColor, color.White))
	canvas.Circle(
		int((float32(x)+float32(3.5))*float32(scale)),
		int((float32(y)+float32(3.5))*float32(scale)),
		int(float32(3.5)*float32(scale)), StrokeStyle(targetColor, color.Black),
	)
	canvas.Circle(
		int((float32(x)+float32(3.5))*float32(scale)),
		int((float32(y)+float32(3.5))*float32(scale)),
		int(float32(2.5)*float32(scale)), StrokeStyle(targetColor, color.White),
	)
	canvas.Circle(
		int((float32(x)+float32(3.5))*float32(scale)),
		int((float32(y)+float32(3.5))*float32(scale)),
		int(float32(1.5)*float32(scale)), StrokeStyle(targetColor, color.Black),
	)
}

func NoStrokeStyle(target, source color.Color) string {
	if source == color.Black {
		return fmt.Sprintf("fill:%s;stroke:none", ColorToFill(target))
	}
	return "fill:white;stroke:none"
}

func StrokeStyle(target, source color.Color) string {
	if source == color.Black {
		return fmt.Sprintf("fill:%s;stroke:white;stroke-width:2", ColorToFill(target))
	}
	return "fill:white;stroke:white;stroke-width:2"
}

func ColorToFill(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("rgb(%d %d %d)", r>>8, g>>8, b>>8)
}
