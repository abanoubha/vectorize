package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
)

const threshold = 30

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
	// convert to SVG
	err = toSVG(img, "test_cases/stat.svg")
	if err != nil {
		log.Fatal(err)
	}

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

func toSVG(img image.Image, svgFilename string) error {
	tempDir := os.TempDir()
	tempPNG, err := os.CreateTemp(tempDir, "temp*.png")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempPNG.Name())
	err = png.Encode(tempPNG, img)
	if err != nil {
		log.Fatal(err)
	}
	tempPNG.Close()

	tempPNM := filepath.Join(tempDir, "temp.pnm")
	defer os.Remove(tempPNM)

	// 3. Encode and write PNM data
	pnmFile, err := os.Create(tempPNM)
	if err != nil {
		log.Fatal(err)
	}
	defer pnmFile.Close()

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Fprintf(pnmFile, "P3\n") // P3 for ASCII color
	fmt.Fprintf(pnmFile, "%d %d\n", width, height)
	fmt.Fprintf(pnmFile, "255\n") // Max color value

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			fmt.Fprintf(pnmFile, "%d %d %d ", r>>8, g>>8, b>>8) // Scale to 0-255
		}
		fmt.Fprintf(pnmFile, "\n")
	}

	// apt install potrace
	cmd := exec.Command("potrace", "-s", "-o", svgFilename, tempPNM) // -s for SVG output

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()
	if err != nil {
		fmt.Println("potrace error:", err)
		fmt.Println("stderr:", errb.String()) // Print stderr
		return err
	}

	fmt.Println("potrace stdout:", outb.String()) // Print stdout
	fmt.Println("PNG converted to SVG successfully!")
	return nil
}
