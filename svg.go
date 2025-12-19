package main

import (
	"fmt"
	"image/color"
	"os"

	svg "github.com/ajstarks/svgo"
)

const path string = "qr.svg"

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

	rows := len(req.Pixels)
	cols := len(req.Pixels[0])
	width := cols * req.Scale
	height := rows * req.Scale
	canvas := svg.New(file)
	canvas.Start(width, height)

	for y, row := range req.Pixels {
		for x, c := range row {
			switch req.Shape {
			case ShapeCircle:
				canvas.Circle(x*req.Scale+req.Scale/2, y*req.Scale+req.Scale/2, req.Scale/2, Style(c))
			case ShapeSquare:
				canvas.Square(x*req.Scale, y*req.Scale, req.Scale, Style(c))
			default:
				fmt.Println("Not implemented shape:", req.Shape)
			}
		}
	}
	canvas.End()
}

func Style(c color.Color) string {
	if c == color.Black {
		return "fill:black;stroke:none"
	}
	return "fill:white;stroke:none"
}

// _basicDot(args: BasicFigureDrawArgs): void {
//   const { size, x, y } = args;
//
//   this._rotateFigure({
//     ...args,
//     draw: () => {
//       this._element = this._window.document.createElementNS("http://www.w3.org/2000/svg", "circle");
//       this._element.setAttribute("cx", String(x + size / 2));
//       this._element.setAttribute("cy", String(y + size / 2));
//       this._element.setAttribute("r", String(size / 2));
//     }
//   });
// }
