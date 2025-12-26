package images

import (
	"image"
	"image/color"
	_ "image/gif"  // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
	"image/png"    // Register PNG encoder
	"math"
	"os"
)

func MakeLogo(path string) string {
	// 1. Open the original image
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 2. Decode the image (handles jpg, png, gif if imported)
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// 3. Create a new image with Alpha channel (RGBA)
	// We need a new canvas because JPEGs generally don't have transparency
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	// 4. Trim image to centered square
	size := min(bounds.Dx(), bounds.Dy())
	offsetX := (bounds.Dx() - size) / 2
	offsetY := (bounds.Dy() - size) / 2
	for y := range size {
		for x := range size {
			originalColor := img.At(x+offsetX, y+offsetY)
			r, g, b, _ := originalColor.RGBA()
			// 5. Set alpha to 0
			rgba.Set(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: 0,
			})
		}
	}

	// 5. Create circular mask
	center := size / 2
	radius := size / 2
	for y := range size {
		for x := range size {
			dx := x - center
			dy := y - center
			distance := math.Sqrt(float64(dx*dx + dy*dy))
			if distance > float64(radius) {
				// Outside the circle, set alpha to 0
				r, g, b, _ := rgba.At(x, y).RGBA()
				rgba.Set(x, y, color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: 0,
				})
			} else {
				// Inside the circle, set alpha to 255
				r, g, b, _ := rgba.At(x, y).RGBA()
				rgba.Set(x, y, color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: 255,
				})
			}
		}
	}

	// 6. Save as PNG
	outFile, err := os.Create("logo.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Encode takes the writer and the image
	err = png.Encode(outFile, rgba)
	if err != nil {
		panic(err)
	}
	return "logo.png"
}
