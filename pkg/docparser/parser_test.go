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

func TestParsePdf(t *testing.T) {

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

func TestParseUrl(t *testing.T) {

	// Arrange
	const url string = "https://hagakure.substack.com/p/twh45-against-overwhelm"

	// Act
	res, err := ParseUrl(url)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("Expected a non-empty string")
	}
	fmt.Println(res)
	fmt.Print(len(res[1]))
}

func TestParseUrlv2(t *testing.T) {

	// Arrange
	const url string = "https://hagakure.substack.com/p/twh45-against-overwhelm"

	// Act
	res, err := getBodyTextFromURLTest(url)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if res == "" {
		t.Error("Expected a non-empty string")
	}
	fmt.Println(res)
	fmt.Print(len(res))
}
