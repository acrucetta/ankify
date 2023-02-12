package docparser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"

	"github.com/joho/godotenv"
)

func ParseTxt(txt_path string) (string, error) {
	// Load environment variables
	err_env := godotenv.Load(".env")
	if err_env != nil {
		fmt.Printf("Error: %v", err_env)
		return "", err_env
	}

	// Load the text file
	f, err := os.Open(txt_path)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}
	defer f.Close()

	// Read the text file
	reader := bufio.NewReader(f)
	var text string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		text += line
	}
	return text, nil
}

func ParsePdf(pdf_path string) (string, error) {

	// Load environment variables
	err_env := godotenv.Load(".env")
	if err_env != nil {
		fmt.Printf("Error: %v", err_env)
		return "", err_env
	}

	// Initialize the library
	err := license.SetMeteredKey(os.Getenv("PDF_PARSER_KEY"))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}
	// Load the PDF document
	reader, f, err := model.NewPdfReaderFromFile(pdf_path, nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}
	defer f.Close()

	// Print number of pages in the PDF.
	numPages, _ := reader.GetNumPages()
	fmt.Printf("Number of pages in the PDF: %d\n", numPages)

	// Ask the user for a page number.
	fmt.Printf("Enter a page number (1-%d): ", numPages)
	user_input := bufio.NewReader(os.Stdin)
	var pageNum string
	pageNum, err = user_input.ReadString('\n')

	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}

	// Convert the page number to an integer.
	// First we trim and remove the newline character.
	pageNum = pageNum[:len(pageNum)-1]
	pageNumInt, err := strconv.Atoi(pageNum)

	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}

	// Get the page.
	page, err := reader.GetPage(pageNumInt)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}

	// Extract text from the page.
	ex, err := extractor.New(page)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}

	text, err := ex.ExtractText()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return "", err
	}
	return text, nil
}
