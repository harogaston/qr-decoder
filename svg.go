package main

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
				canvas.Circle(x*req.Scale+req.Scale/2, y*req.Scale+req.Scale/2, req.Scale/2, CircleStyle(c))
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
	canvas.End()
}

func SquareStyle(c color.Color) string {
	if c == color.Black {
		return "fill:black;stroke:none"
	}
	return "fill:white;stroke:none"
}

func CircleStyle(c color.Color) string {
	if c == color.Black {
		return "fill:black;stroke:white;stroke-width:2"
	}
	return "fill:white;stroke:white;stroke-width:2"
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
