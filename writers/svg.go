package writers

import (
	"fmt"
	"image/color"
	"os"

	svg "github.com/twpayne/go-svg"
	"github.com/twpayne/go-svg/svgpath"
)

const output_file_path string = "qr.svg"
const logoRelativeSize = 1. / 5.

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

	// Draw logo ensuring a minimum size of 5 modules
	logoSize := int(float64(dim) * logoRelativeSize)
	logoSize += (logoSize ^ dim) & 1
	fmt.Println("Logo size:", logoSize)

	if req.Logo != "" && logoSize >= 5 {
		logoPos := dim/2 - logoSize/2

		// Cleanup overlapping QR modules
		padding := 0.
		switch req.Shape {
		case ShapeCircle:
			padding = 1
		case ShapeSquircle:
			padding = 2
		case ShapeSquare:
			padding = 0
		}
		startCell := logoPos
		endCell := startCell + logoSize

		center := float64(logoPos) + float64(logoSize)/2.
		radius := float64(logoSize/2) + padding
		for y := startCell; y < endCell; y++ {
			for x := startCell; x < endCell; x++ {
				dx := float64(x) + .5 - center
				dy := float64(y) + .5 - center
				distance := dx*dx + dy*dy
				if distance < radius*radius {
					canvas.AppendChildren(
						svg.Use().XY(float64(x), float64(y), svg.Number).Href("#square").Style("fill:red"),
					)
				}
			}
		}

		// Place logo with clipping path
		logoClipPath := svg.Use().Href(svg.String("#" + string(req.Shape)))
		logoClipPath.Attrs["transform"] = svg.String(fmt.Sprintf("scale(%d) translate(%f %f)", logoSize, float64(logoPos)/float64(logoSize), float64(logoPos)/float64(logoSize)))
		canvas.AppendChildren(
			svg.ClipPath().ID("logoClip").AppendChildren(
				logoClipPath,
			),
			svg.Image().Href(svg.String(req.Logo)).XYWidthHeight(
				float64(logoPos), float64(logoPos), float64(logoSize), float64(logoSize), svg.Number,
			).ClipPath("url(#logoClip)"),
		)
	}
	if _, err := canvas.WriteToIndent(file, "", "  "); err != nil {
		panic(err)
	}
}

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
