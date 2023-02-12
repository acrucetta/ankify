package docparser

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func ParseTxt(txt_path string) (map[int]string, error) {

	// Load the text file
	f, err := os.Open(txt_path)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
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
	res := make(map[int]string)
	res[1] = text

	return res, nil
}

func GetMaxPages(pdf_path string) (int, error) {
	// Read the PDF file
	config := pdfcpu.NewDefaultConfiguration()
	err := api.ExtractMetadataFile(pdf_path, "output", config)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return 0, err
	}
	return 0, nil
}

// This function is used to parse a PDF file and return the text on a specific page.
func ParsePdfPage(pdf_path string, page_number int) (string, error) {

	// Print current path
	dir, _ := os.Getwd()
	fmt.Println(dir)

	// Call pdf2text.py to parse the PDF file on the given page
	command := "pipenv run pdf2txt.py -p " + strconv.Itoa(page_number) + " " + pdf_path
	cmd := exec.Command("bash", "-c", command)

	// Run the command
	out, err := cmd.Output()

	// Check for errors
	if err != nil {
		fmt.Printf("Error: %v", err)
		return "", err
	}
	return string(out), nil
}

// This function is used to parse a PDF file and return the text on a specific page.
// We are using pdf2text.py to parse the PDF file. We will call it from the command line.
// The function returns the text on the page as a string.
func ParsePdf(pdf_path string, pages []int) (map[int]string, error) {

	// Check the pages that were requested are valid
	for _, page := range pages {
		if page < 1 {
			return nil, fmt.Errorf("Invalid page number: %v", page)
		}
	}

	// Call pdf2text.py to parse the PDF file on each page
	// and store the output in a dictionary.
	var parsed_pages map[int]string = make(map[int]string)

	for _, page := range pages {
		parsed_page, err := ParsePdfPage(pdf_path, page)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return nil, err
		}
		parsed_pages[page] = parsed_page
	}

	return parsed_pages, nil
}
