package main

import "fmt"

func main() {
	// Ask the user for their API key
	var API_KEY string
	fmt.Print("Enter your OpenAI API key: ")
	fmt.Scanln(&API_KEY)

	// Ask the user for their pdf path
	var PDF_PATH string
	fmt.Print("Enter the path to your pdf: ")
	fmt.Scanln(&PDF_PATH)

}
