package docparser

import (
	"testing"
)

func TestParsePdf(t *testing.T) {

	// Arrange
	const pdfPath string = "../../data/test.pdf"

	// Act
	res, err := ParsePdf(pdfPath)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if res == "" {
		t.Error("Expected a non-empty string")
	}
}
