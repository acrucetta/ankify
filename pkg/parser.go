package pdfparser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
)

func parse_pdf(pdfPath string) (string, error) {

	// Initialize the library
	err := license.SetMeteredKey(os.Getenv("PDF_PARSER_KEY"))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}
	// Load the PDF document
	reader, f, err := model.NewPdfReaderFromFile(pdfPath, nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}
	defer f.Close()

	// Print number of pages in the PDF.
	numPages, _ := reader.GetNumPages()
	fmt.Printf("Number of pages in the PDF: %d\n", numPages)

	// Ask the user for a page number.
	fmt.Printf("Enter a page number: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	pageNum, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}

	// Get the page.
	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}

	// Extract text from the page.
	text, err := page.GetText()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return "", err
	}

	// Print the text.
	fmt.Println(text)
}
