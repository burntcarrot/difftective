package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/gen2brain/go-fitz"
)

// CompareFiles loads images from files and then compares it with CompareImages.
func CompareFiles(src, dst string) (diff image.Image, percent float64, err error) {
	srcImage, err := loadImage(src)
	if err != nil {
		return nil, 0, err
	}
	dstImage, err := loadImage(dst)
	if err != nil {
		return nil, 0, err
	}
	return CompareImages(srcImage, dstImage)
}

// CompareImages checks image size and returns the difference after pixel-by-pixel comparison.
func CompareImages(src, dst image.Image) (diff image.Image, percent float64, err error) {
	// Check if the images have the same dimensions.
	srcBounds := src.Bounds()
	dstBounds := dst.Bounds()
	if !boundsMatch(srcBounds, dstBounds) {
		return nil, 100.0, errors.New("image sizes don't match")
	}

	// Create a new difference image.
	diffImage := image.NewRGBA(image.Rect(0, 0, srcBounds.Max.X, srcBounds.Max.Y))

	// Perform pixel-by-pixel matching.
	var differentPixels float64
	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
			r, g, b, _ := dst.At(x, y).RGBA()

			// Set the pixel regardless of the difference.
			diffImage.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})

			// If there exists some amount of difference, proceed with coloring the different pixel with a red color.
			if !isEqualColor(src.At(x, y), dst.At(x, y)) {
				differentPixels++

				// Use red color for displaying difference.
				diffImage.Set(x, y, color.RGBA{255, 0, 0, 255})
			}
		}
	}

	// Calculate the difference percentage.
	diffPercent := differentPixels / float64(srcBounds.Max.X*srcBounds.Max.Y) * 100

	return diffImage, diffPercent, nil
}

// isEqualColor checks if two colors are equal or not.
// Uses the RGBA colorspace.
func isEqualColor(a, b color.Color) bool {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()

	// Should be an exact match. If not a match, it returns false.
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// boundsMatch checks if the dimensions of the images are equal or not.
func boundsMatch(a, b image.Rectangle) bool {
	return a.Min.X == b.Min.X && a.Min.Y == b.Min.Y && a.Max.X == b.Max.X && a.Max.Y == b.Max.Y
}

// loadImage reads an image from the filesystem, and returns an image.
func loadImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Decode image from the file reader.
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// writeImage writes an image to a file on the user's filesystem.
func writeImage(filename string, img image.Image) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	// Encode image in PNG format.
	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	defer out.Close()

	return nil
}

// pdfToImage reads the PDF and converts it to an image.
// Makes use of the MuPDF fitz driver.
func pdfToImage(filename, imgName string) error {
	doc, err := fitz.New(filename)
	if err != nil {
		return err
	}

	defer doc.Close()

	pageNumber := 0
	img, err := doc.Image(pageNumber)
	if err != nil {
		return err
	}

	err = writeImage(imgName, img)
	if err != nil {
		return err
	}

	return nil
}
