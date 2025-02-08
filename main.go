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
	imgFile, err := os.Open("test_cases/stat.png")
	if err != nil {
		log.Fatal(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	// processing the image
	// use less colors
	img = useLessColors(img)
	// TODO: convert to SVG

	// Save processed image directly
	outFile, err := os.Create("test_cases/stat-processed.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img) // Use original image for saving
	if err != nil {
		log.Fatal(err)
	}
}

func useLessColors(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Define color comparison function
	compareColor := func(c1, c2 color.Color) bool {
		r1, g1, b1, _ := c1.RGBA()
		r2, g2, b2, _ := c2.RGBA()
		return math.Abs(float64(r1-r2)) < threshold &&
			math.Abs(float64(g1-g2)) < threshold &&
			math.Abs(float64(b1-b2)) < threshold
	}

	// Define color averaging function
	averageColor := func(colors []color.Color) color.Color {
		var r, g, b, a uint32
		for _, c := range colors {
			cr, cg, cb, ca := c.RGBA()
			r += cr
			g += cg
			b += cb
			a += ca
		}
		n := uint32(len(colors))
		return color.RGBA64{uint16(r / n), uint16(g / n), uint16(b / n), uint16(a / n)}
	}

	// Process image
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			p := img.At(x, y)
			neighbors := []color.Color{p}

			// Check horizontal neighbors
			if x > 0 && compareColor(p, img.At(x-1, y)) {
				neighbors = append(neighbors, img.At(x-1, y))
			}

			// Check vertical neighbors
			if y > 0 && compareColor(p, img.At(x, y-1)) {
				neighbors = append(neighbors, img.At(x, y-1))
			}

			// Set the pixel to the average color of similar neighbors
			newImg.Set(x, y, averageColor(neighbors))
		}
	}

	return newImg
}
