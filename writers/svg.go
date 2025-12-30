package writers

import (
	"fmt"
	"image/color"
	"math/rand/v2"

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
						svg.String(GetStyle(req.Shape, req.Color, c)),
					),
				)
			}
		}
	}

	// Connect neighboring modules
	connect(req, canvas, dim)

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
			padding = 2
		case ShapeSquircle:
			padding = 3
		case ShapeSquare:
			padding = 2
		}
		startCell := logoPos
		endCell := startCell + logoSize

		center := float64(logoPos) + float64(logoSize)/2.
		radius := float64(logoSize/2) + padding
		for y := startCell - 1; y < endCell+1; y++ {
			for x := startCell - 1; x < endCell+1; x++ {
				dx := float64(x) + .5 - center
				dy := float64(y) + .5 - center
				distance := dx*dx + dy*dy
				if distance < radius*radius {
					canvas.AppendChildren(
						svg.Use().XY(float64(x), float64(y), svg.Number).Href("#square").Style("fill:white"),
					)
				}
			}
		}

		// Place logo with clipping path
		logoClipPath := svg.Use().Href("#circle")
		logoClipPath.Attrs["transform"] = svg.String(fmt.Sprintf("scale(%d) translate(%f %f)", logoSize, float64(logoPos)/float64(logoSize), float64(logoPos)/float64(logoSize)))

		logoBorderScale := float64(logoSize) + .8
		logoBorderPos := float64(dim)/2. - float64(logoBorderScale)/2.
		fmt.Println("Logo border position:", logoBorderPos)
		fmt.Println("Logo border scale:", logoBorderScale)
		logoBorder := svg.Use().Href("#circle").Style(svg.String(fmt.Sprintf("fill:none;stroke:black;stroke-width:%f", 0.25/float64(logoBorderScale))))
		logoBorder.Attrs["transform"] = svg.String(fmt.Sprintf("scale(%f) translate(%f %f)", logoBorderScale, float64(logoBorderPos)/float64(logoBorderScale), float64(logoBorderPos)/float64(logoBorderScale)))
		canvas.AppendChildren(
			svg.ClipPath().ID("logoClip").AppendChildren(
				logoClipPath,
			),
			logoBorder,
			svg.Image().Href(svg.String(req.Logo)).XYWidthHeight(
				float64(logoPos), float64(logoPos), float64(logoSize), float64(logoSize), svg.Number,
			).ClipPath("url(#logoClip)"),
		)
	}
	if _, err := canvas.WriteToIndent(file, "", "  "); err != nil {
		panic(err)
	}
}

// connect encapsulates the logic for drawing connected shapes (rectangles) based on module color.
func connect(req SVGRequest, canvas *svg.SVGElement, dim int) {
	// featEnabled is false in the original code, thus this whole function is effectively disabled.
	// We preserve the original behavior.
	featEnabled := false
	if !featEnabled {
		return
	}

	width := len(req.Cells) // Using len(req.Cells) for correct indexing into the visited map
	visited := make(map[int]struct{}, width*width)

	// Helper function to draw a connecting rectangle
	drawConnectingRect := func(x, y int, dir uint64, c color.Color) {
		if dir == 0 { // Horizontal
			canvas.AppendChildren(
				svg.Rect().XYWidthHeight(float64(x-4)+0.5, float64(y-4), 1, 1, svg.Number).Style(
					svg.String(NoStrokeStyle(req.Color, c)),
				),
			)
		} else { // Vertical
			canvas.AppendChildren(
				svg.Rect().XYWidthHeight(float64(x-4), float64(y-4)+0.5, 1, 1, svg.Number).Style(
					svg.String(NoStrokeStyle(req.Color, c)),
				),
			)
		}
	}

	// Helper function to attempt connection in a given direction
	tryConnectInDirection := func(startX, startY int, c color.Color, currentDir uint64) bool {
		foundNeighbor := false
		if currentDir == 0 { // Horizontal
			for nx := startX + 1; nx < dim && nx < width && req.Cells[startY][nx] == c; nx++ { // Original code used `nx < dim`. Preserving it.
				if _, ok := visited[startY*width+nx]; ok {
					break
				}
				foundNeighbor = true
				visited[startY*width+nx] = struct{}{}
				drawConnectingRect(nx-1, startY, currentDir, c) // nx-1 is the x-coordinate of the left module for the connection
			}
		} else { // Vertical
			for ny := startY + 1; ny < width && req.Cells[ny][startX] == c; ny++ {
				if _, ok := visited[ny*width+startX]; ok {
					break
				}
				foundNeighbor = true
				visited[ny*width+startX] = struct{}{}
				drawConnectingRect(startX, ny-1, currentDir, c) // ny-1 is the y-coordinate of the top module for the connection
			}
		}
		return foundNeighbor
	}

	for y, row := range req.Cells {
		for x, c := range row {
			if c == color.Black { // Only consider black modules
				if _, ok := visited[y*width+x]; !ok {
					visited[y*width+x] = struct{}{} // Mark current cell as visited

					dir := rand.Uint64() & 1 // Pick random initial direction (0 for horizontal, 1 for vertical)

					// Try connecting in the first direction
					if !tryConnectInDirection(x, y, c, dir) {
						// If no connection was found in the first direction, try the other direction
						tryConnectInDirection(x, y, c, 1-dir)
					}
				}
			}
		}
	}
}

func GetStyle(shape Shape, target, source color.Color) string {
	switch shape {
	case ShapeSquare:
		return NoStrokeStyle(target, source)
	default:
		return StrokeStyle(target, source)
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
