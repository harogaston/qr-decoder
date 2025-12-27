package writers

import (
	"fmt"
	"image/color"
	"os"

	// "github.com/harogaston/qr-decoder/images"
	svg "github.com/twpayne/go-svg"
	"github.com/twpayne/go-svg/svgpath"
)

const output_file_path string = "qr.svg"
const roundedCornerPath = "M${x},${y} v ${size} h ${size} v ${-size/2} a ${size/2},${size/2} 0,0,0 ${-size/2},${-size/2} z"
const squircleCurvature = 0.75
const squirclePath = "M0,${size/2} C 0 ${curvA}, ${curvA} 0, ${size/2} 0 S ${size} ${curvA}, ${size} ${size/2}, ${curvB} ${size}, ${size/2} ${size}, 0 ${curvB}, 0 ${size/2}"

const logoRelativeSize = 5

type SVGRequest struct {
	Scale int
	Cells [][]color.Color
	Shape Shape
	Logo  string
	Color color.Color
}

func WriteSVG(req SVGRequest) {
	if req.Color == nil {
		req.Color = color.Black
	}

	file, err := os.Create(output_file_path)
	if err != nil {
		fmt.Println("Error creating SVG file:", err)
	}
	defer file.Close()

	dim := len(req.Cells) - 8
	canvasSize := float64(dim + 4)
	canvas := svg.New().
		WidthHeight(canvasSize, canvasSize, svg.Number).
		Transform(svg.String(fmt.Sprintf("scale(%d) translate(4, 4)", req.Scale)))
	canvas.Attrs["transform-origin"] = svg.String("0 0")

	// Definitions
	circle := svg.Circle().R(svg.Number(0.5)).ID(svg.String(ShapeCircle))
	circle.Attrs["transform"] = svg.String("translate(0.5,0.5)")
	canvas.AppendChildren(
		svg.Defs(
			circle,
			svg.Rect().XYWidthHeight(0, 0, 1, 1, svg.Number).ID(svg.String(ShapeSquare)),
			svg.Path().D(
				svgpath.New().MoveToAbs([]float64{0, 0.5}).CurveToAbs([]float64{0, 0.125}, []float64{0.125, 0}, []float64{0.5, 0}).SCurveToAbs([]float64{1, 0.125}, []float64{1., 0.5}, []float64{0.875, 1.}, []float64{0.5, 1.}, []float64{0, 0.875}, []float64{0, 0.5}).ClosePath(),
			).ID(svg.String(ShapeSquircle)),
		),
	)

	// Finder pattern block
	// canvas.Group(`id="finderpattern"`)
	// CircleFinderPattern(canvas, req.Color, 4, 4, req.Scale)
	// CircleFinderPattern(canvas, req.Color, (dim - 7 - 4), 4, req.Scale)
	// CircleFinderPattern(canvas, req.Color, 4, (dim - 7 - 4), req.Scale)
	// canvas.Gend()
	// canvas.Square(0, 0, req.Scale, `id="cell"`)
	// squirclePathExp := strings.ReplaceAll(squirclePath, "${size}", fmt.Sprintf("%d", req.Scale))
	// squirclePathExp = strings.ReplaceAll(squirclePathExp, "${size/2}", fmt.Sprintf("%d", req.Scale/2))
	// squirclePathExp = strings.ReplaceAll(squirclePathExp, "${curvA}", fmt.Sprintf("%.f", float32(req.Scale/2.)*float32(1.-squircleCurvature)))
	// squirclePathExp = strings.ReplaceAll(squirclePathExp, "${curvB}", fmt.Sprintf("%.f", float32(req.Scale)-float32(req.Scale/2.)*float32(1.-squircleCurvature)))
	// canvas.Path(squirclePathExp, `id="cell"`, `transform="rotate(0,80,80)"`)
	//
	// // Finder pattern block
	// canvas.Group(`id="finderpattern"`)
	// SquircleFinderPattern(canvas, req.Color, 4, 4, req.Scale)
	// SquircleFinderPattern(canvas, req.Color, (dim - 7 - 4), 4, req.Scale)
	// SquircleFinderPattern(canvas, req.Color, 4, (dim - 7 - 4), req.Scale)
	// canvas.Gend()
	// pathExpression := strings.ReplaceAll(roundedCornerPath, "${x}", fmt.Sprintf("%d", x*req.Scale))
	// pathExpression = strings.ReplaceAll(pathExpression, "${y}", fmt.Sprintf("%d", y*req.Scale))
	// pathExpression = strings.ReplaceAll(pathExpression, "${size}", fmt.Sprintf("%d", req.Scale))
	// pathExpression = strings.ReplaceAll(pathExpression, "${-size/2}", fmt.Sprintf("%d", -req.Scale/2))
	// pathExpression = strings.ReplaceAll(pathExpression, "${size/2}", fmt.Sprintf("%d", req.Scale/2))
	// canvas.Path(pathExpression, SquareStyle(c))

	// All Modules
	for y, row := range req.Cells {
		for x, c := range row {
			if c == color.Black {
				canvas.AppendChildren(
					svg.Use().XY(float64(x-4), float64(y-4), svg.Number).Href(svg.String(fmt.Sprintf("#%s", req.Shape))).Style(
						svg.String(StrokeStyle(req.Color, c)),
					),
				)
			}
		}
	}

	// Superimpose finder patterns
	finderBackground := svg.Use().Href(svg.String("#square")).Style("fill:white")
	finderBackground.Attrs["transform"] = svg.String(fmt.Sprintf("scale(%d) translate(%f, %f)", 7, 0/7., 0/7.))
	finderOuterRing := svg.Use().Href(svg.String("#" + string(req.Shape))).Style(svg.String(NoStrokeStyle(req.Color, color.Black)))
	finderOuterRing.Attrs["transform"] = svg.String(fmt.Sprintf("scale(%d) translate(%f, %f)", 7, 0/7., 0/7.))
	finderMiddleRing := svg.Use().Href(svg.String("#" + string(req.Shape))).Style(svg.String(NoStrokeStyle(req.Color, color.White)))
	finderMiddleRing.Attrs["transform"] = svg.String(fmt.Sprintf("scale(%d) translate(%f, %f)", 5, 1./5., 1./5.))
	finderCenterRing := svg.Use().Href(svg.String("#" + string(req.Shape))).Style(svg.String(NoStrokeStyle(req.Color, color.Black)))
	finderCenterRing.Attrs["transform"] = svg.String(fmt.Sprintf("scale(%d) translate(%f, %f)", 3, 2./3., 2./3.))
	canvas.AppendChildren(
		svg.G(
			finderBackground,
			finderOuterRing,
			finderMiddleRing,
			finderCenterRing,
		).ID("finderpattern"),
		svg.Use().Href(svg.String("#finderpattern")).XY(float64(dim-7), 0, svg.Number),
		svg.Use().Href(svg.String("#finderpattern")).XY(0, float64(dim-7), svg.Number),
	)

	// Superimpose logo
	if req.Logo != "" {
		logoClipPath := svg.Use().Href(svg.String("#" + string(req.Shape)))
		logoClipPath.Attrs["transform"] = svg.String(fmt.Sprintf("scale(%f)", float64(dim/5)))
		canvas.AppendChildren(
			svg.ClipPath(logoClipPath).ID("logoClipPath"),
			svg.Image().ClipPath(svg.String("url(#logoClipPath)")).Href(svg.String(req.Logo)).XYWidthHeight(
				float64(dim)/2., float64(dim)/2., float64(dim)/5., float64(dim)/5., svg.Number,
			),
		)
		// 	// Create rounded logo version
		// 	logoPath := images.MakeLogo(req.Logo)
		// 	logoSize := dim / 5 * req.Scale
		//
		// 	// Delete overlapping QR modules
		// 	// For all cells inside the square of logoSize, if any cell corners overlap, draw white square
		// 	logoDim := dim/logoRelativeSize + 1 // logo size in cells
		//
		// 	startCell := (dim - logoDim) / 2
		// 	endCell := startCell + logoDim
		//
		// 	for y := startCell; y <= endCell; y++ {
		// 		for x := startCell; x <= endCell; x++ {
		// 			// Check if inside logo circle area
		// 			center := dim / 2
		// 			radius := logoDim / 2
		// 			radius += 1 // Add some padding
		// 			dx := x - center
		// 			dy := y - center
		// 			distance := dx*dx + dy*dy
		// 			if distance < radius*radius {
		// 				canvas.Square(x*req.Scale, y*req.Scale, req.Scale, NoStrokeStyle(req.Color, color.White))
		// 			}
		// 		}
		// 	}
		//
		// 	// Place logo at center
		// 	logoPos := (canvasSize - logoSize) / 2
		// 	canvas.Image(logoPos, logoPos, logoSize, logoSize, logoPath)
	}
	if _, err := canvas.WriteToIndent(file, "", "  "); err != nil {
		panic(err)
	}
}

// func SquircleFinderPattern(canvas *svg.SVG, targetColor color.Color, x, y int, scale int) {
// 	canvas.Square(x*scale, y*scale, 8*scale, NoStrokeStyle(targetColor, color.White))
//
// 	squirclePathExp := strings.ReplaceAll(squirclePath, "${size}", fmt.Sprintf("%d", 8*scale))
// 	squirclePathExp = strings.ReplaceAll(squirclePathExp, "${size/2}", fmt.Sprintf("%d", 8*scale/2))
// 	squirclePathExp = strings.ReplaceAll(squirclePathExp, "${curvA}", fmt.Sprintf("%.f", float32(8*scale/2.)*float32(1.-squircleCurvature)))
// 	squirclePathExp = strings.ReplaceAll(squirclePathExp, "${curvB}", fmt.Sprintf("%.f", float32(8*scale)-float32(scale/2.)*float32(1.-squircleCurvature)))
// 	canvas.Path(squirclePathExp, fmt.Sprintf(`transform="rotate(0,80,80) scale(%d)"`, 1), fmt.Sprintf(`transform-origin="%d %d"`, x*scale, y*scale))
// 	// canvas.Use(x*scale, y*scale, "#cell", StrokeStyle(targetColor, color.Black), fmt.Sprintf(`transform="scale(%d)"`, 8*scale))
// 	// canvas.Use((x+1)*scale, (y+1)*scale, "#cell", StrokeStyle(targetColor, color.White), fmt.Sprintf(`transform="scale(%d)"`, 6*scale))
// 	// canvas.Use((x+2)*scale, (y+2)*scale, "#cell", StrokeStyle(targetColor, color.Black), fmt.Sprintf(`transform="scale(%d)"`, 4*scale))
// }
//
// func SquareFinderPattern(canvas *svg.SVG, targetColor color.Color, x, y int, scale int) {
// 	canvas.Square(x*scale, y*scale, 8*scale, NoStrokeStyle(targetColor, color.White))
// 	// canvas.Square((x)*scale, (y)*scale, 7*scale, StrokeStyle(targetColor, color.Black))
// 	// canvas.Square((x+1)*scale, (y+1)*scale, 5*scale, StrokeStyle(targetColor, color.White))
// 	canvas.Use((x+1)*scale, (y+1)*scale, "#cell", fmt.Sprintf(`transform-origin="%d %d"`, (x+1)*scale, (y+1)*scale), `transform="scale(5)"`, StrokeStyle(targetColor, color.Black))
// 	canvas.Use((x+2)*scale, (y+2)*scale, "#cell", fmt.Sprintf(`transform-origin="%d %d"`, (x+2)*scale, (y+2)*scale), `transform="scale(3)"`, StrokeStyle(targetColor, color.White))
//
// 	// canvas.Square((x+2)*scale, (y+2)*scale, 3*scale, StrokeStyle(targetColor, color.Black))
// }
//
// func CircleFinderPattern(canvas *svg.SVG, targetColor color.Color, x, y int, scale int) {
// 	canvas.Square(x*scale, y*scale, 8*scale, NoStrokeStyle(targetColor, color.White))
// 	canvas.Circle(
// 		int((float32(x)+float32(3.5))*float32(scale)),
// 		int((float32(y)+float32(3.5))*float32(scale)),
// 		int(float32(3.5)*float32(scale)), StrokeStyle(targetColor, color.Black),
// 	)
// 	canvas.Circle(
// 		int((float32(x)+float32(3.5))*float32(scale)),
// 		int((float32(y)+float32(3.5))*float32(scale)),
// 		int(float32(2.5)*float32(scale)), StrokeStyle(targetColor, color.White),
// 	)
// 	canvas.Circle(
// 		int((float32(x)+float32(3.5))*float32(scale)),
// 		int((float32(y)+float32(3.5))*float32(scale)),
// 		int(float32(1.5)*float32(scale)), StrokeStyle(targetColor, color.Black),
// 	)
// }

func NoStrokeStyle(target, source color.Color) string {
	if source == color.Black {
		return fmt.Sprintf("fill:%s;stroke:none", ColorToFill(target))
	}
	return "fill:white;stroke:none"
}

func StrokeStyle(target, source color.Color) string {
	if source == color.Black {
		return fmt.Sprintf("fill:%s;stroke:white;stroke-width:0.125", ColorToFill(target))
	}
	return "fill:white;stroke:white;stroke-width:0.125"
}

func ColorToFill(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("rgb(%d %d %d)", r>>8, g>>8, b>>8)
}
