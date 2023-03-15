package docparser

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/net/html"
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

// This function is used to extract the raw text from a URL.
func ParseUrl(url string) (map[int]string, error) {
	// Ensure that the URL starts with http:// or https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	// Check if the URL is valid, return loggin error if not
	_, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	bodyText, err := getBodyTextFromURL(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var body_text map[int]string = make(map[int]string)
	body_text[1] = bodyText

	return body_text, nil
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

func getBodyTextFromURL(url string) (string, error) {
	// Fetch the HTML document from the URL.
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Parse the HTML document.
	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	// Extract the text from the content elements in the HTML document.
	var text string
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "h1", "h2", "h3", "h4", "h5", "h6", "li", "em", "ul", "ol", "pre", "td", "br",
			"tbody":
				text += " " + extractNodeText(n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)

	return text, nil
}

func extractNodeText(n *html.Node) string {
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
		}
	}
	return text
}
