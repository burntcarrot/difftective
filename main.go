package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/burntcarrot/ricecake"
)

func main() {
	cli := ricecake.NewCLI("difftective", "Detect differences in images and PDFs", "v0.1.0")

	pdf := cli.NewSubCommand("pdf", "Use PDFs")

	var previousFile string
	var newFile string
	outputFile := "diff.png"

	pdf.StringFlagP("previous", "p", "Path for storing diff output", &previousFile)
	pdf.StringFlagP("new", "n", "Path for storing diff output", &newFile)
	pdf.StringFlagP("output", "o", "Path for storing diff output", &outputFile)

	pdf.Action(func() error {
		oldImage, newImage, err := prepareImages(previousFile, newFile)
		if err != nil {
			return err
		}

		percent, err := generateDiff(oldImage, newImage, outputFile)
		if err != nil {
			return err
		}

		fmt.Printf("Difference: %.3f%%\n", percent)
		fmt.Printf("Saved difference to %s\n", outputFile)

		return nil
	})

	image := cli.NewSubCommand("image", "Use images")

	image.StringFlagP("previous", "p", "Path for storing diff output", &previousFile)
	image.StringFlagP("new", "n", "Path for storing diff output", &newFile)
	image.StringFlagP("output", "o", "Path for storing diff output", &outputFile)

	image.Action(func() error {
		percent, err := generateDiff(previousFile, newFile, outputFile)
		if err != nil {
			return err
		}

		fmt.Printf("Difference: %.3f%%\n", percent)
		fmt.Printf("Saved difference to %s\n", outputFile)

		return nil
	})

	err := cli.Run()
	if err != nil {
		log.Fatalf("failed to run difftective, err: %v", err)
	}
}

// prepareImages prepares images from PDF files.
// Since pixel-by-pixel comparison is performed on images, conversion of a PDF file to an image is required.
func prepareImages(oldPdf, newPdf string) (string, string, error) {
	// The temporary file is used as an intermediate while working with PDF files, where the PDF file is first converted to an image (temporary file).
	// This allows users to use the program without creating a mess in their local filesystem.
	oldPdfImage := getTempFilename(oldPdf)
	err := pdfToImage(oldPdf, oldPdfImage)
	if err != nil {
		return "", "", err
	}

	newPdfImage := getTempFilename(newPdf)
	err = pdfToImage(newPdf, newPdfImage)
	if err != nil {
		return "", "", err
	}

	return oldPdfImage, newPdfImage, err
}

// getTempFilename returns the filename (path) for a temporary file.
func getTempFilename(src string) string {
	tmp := os.TempDir()
	base := filepath.Base(src)

	// Trim extension. difftective currently uses png as output.
	filename := fmt.Sprintf("%s.png", strings.TrimSuffix(base, path.Ext(base)))

	// Construct filename based on the temporary directory.
	filePath := filepath.Join(tmp, filename)

	return filePath
}

// generateDiff compares the old and new image and returns difference.
// Returns -1.0 when it encounters an error.
func generateDiff(oldImage, newImage, dst string) (float64, error) {
	diff, percent, err := CompareFiles(oldImage, newImage)
	if err != nil {
		return -1.0, err
	}

	err = writeImage(dst, diff)
	if err != nil {
		return -1.0, err
	}

	return percent, nil
}
