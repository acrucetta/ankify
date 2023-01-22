package pdfParsing

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	// Load the PDF file from the arguments.
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_parser.go input.pdf\n")
		os.Exit(1)
	}
	pdfPath := os.Args[1]

	// Initialize the library
	err := license.SetMeteredKey(os.Getenv("PDF_PARSER_KEY"))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	// Load the PDF document
	pdfReader, err := model.NewPdfReaderFromFile(pdfPath)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

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
		return
	}

	// Get the page.
	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// Extract text from the page.
	text, err := page.GetText()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Make an API call to OpenAI to generate a question and answer.

}
