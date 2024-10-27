package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

const threshold = 60

func main() {
	imgFile, err := os.Open("test_cases/monkey.png")
	if err != nil {
		log.Fatal(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	// Define color comparison function
	compareColor := func(c1, c2 color.Color) bool {
		r1, g1, b1, _ := c1.RGBA()
		r2, g2, b2, _ := c2.RGBA()
		return math.Abs(float64(r1-r2)) < threshold &&
			math.Abs(float64(g1-g2)) < threshold &&
			math.Abs(float64(b1-b2)) < threshold
	}

	// Process rows directly on the original image
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			p := img.At(x, y)

			// Check horizontal neighbors
			if x > 0 && compareColor(p, img.At(x-1, y)) {
				continue // Skip if color matches left neighbor
			}

			// Check vertical neighbors
			if y > 0 && compareColor(p, img.At(x, y-1)) {
				continue // Skip if color matches above neighbor
			}

			// If no matching neighbors, keep current color (no modification needed)
		}
	}

	// Save processed image directly
	outFile, err := os.Create("test_cases/monkey-processed.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img) // Use original image for saving
	if err != nil {
		log.Fatal(err)
	}
}
