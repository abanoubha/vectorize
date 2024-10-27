package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

func findCommonColors(img image.Image, threshold float64) {
	// Get dimensions of the input image
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a new RGBA image with the same dimensions
	r := image.NewRGBA(image.Rect(0, 0, width, height))

	// Copy pixels manually
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			r.Set(x, y, c)
		}
	}

	// Define color comparison function
	compareColor := func(c1, c2 color.Color) bool {
		r1, g1, b1, _ := c1.RGBA()
		r2, g2, b2, _ := c2.RGBA()
		return math.Abs(float64(r1-r2)) < threshold &&
			math.Abs(float64(g1-g2)) < threshold &&
			math.Abs(float64(b1-b2)) < threshold
	}

	// Process rows
	for y := 0; y < r.Bounds().Dy(); y++ {
		for x := 0; x < r.Bounds().Dx(); x++ {
			p := r.At(x, y)

			// Check horizontal neighbors
			if x > 0 && compareColor(p, r.At(x-1, y)) {
				continue // Skip if color matches left neighbor
			}

			// Check vertical neighbors
			if y > 0 && compareColor(p, r.At(x, y-1)) {
				continue // Skip if color matches above neighbor
			}

			// If no matching neighbors, keep current color
			continue
		}
	}

	// Recreate image with processed colors
	newImg := image.NewRGBA(r.Bounds())

	// Copy pixels from r to newImg
	for y := 0; y < r.Bounds().Dy(); y++ {
		for x := 0; x < r.Bounds().Dx(); x++ {
			newImg.Set(x, y, r.At(x, y))
		}
	}
}
func main() {
	imgFile, err := os.Open("test_cases/nursery-cover.png")
	if err != nil {
		log.Fatal(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	findCommonColors(img, 90) // Adjust threshold as needed

	// Save processed image
	outFile, err := os.Create("test_cases/nursery-cover-processed.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		log.Fatal(err)
	}
}
