package docparser

import (
	"fmt"
	"os"
	"testing"
)

func TestParsePdfPage(t *testing.T) {

	// Print current working directory
	dir, _ := os.Getwd()
	fmt.Println(dir)

	// Arrange
	const pdfPath string = "../../data/test.pdf"

	// Act
	res, err := ParsePdfPage(pdfPath, 1)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if res == "" {
		t.Error("Expected a non-empty string")
	}
}

func TestParsePdfFree(t *testing.T) {

	// Arrange
	const pdfPath string = "../../data/test.pdf"

	// Act
	res, err := ParsePdf(pdfPath, []int{1, 2})

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(res) == 0 {
		t.Error("Expected a non-empty array")
	}
}
