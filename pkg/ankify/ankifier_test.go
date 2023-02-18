package ankify

import (
	"fmt"
	"testing"

	"github.com/acrucetta/anki-builder/pkg/docparser"
	"github.com/joho/godotenv"
)

func TestAnkify(t *testing.T) {

	// Arrange
	const anki_text string = "The capital of France is Paris. The capital of Germany is Berlin. The capital of Spain is Madrid."
	var anki_text_map = make(map[int]string)
	anki_text_map[1] = anki_text

	godotenv.Load("../../.env")

	// Act
	res, err := Ankify(anki_text_map, 10)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if res.Questions == nil {
		t.Error("Expected a non-empty array")
	}

	fmt.Println(res)
}

func TestParseAnkiText(t *testing.T) {

	// Arrange
	const anki_text string = "Q: This is a test\nA: This is a test; Q: This is a test\nA: This is a test"

	// Act
	res, err := ParseAnkiText(anki_text)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(res.Questions) == 0 {
		t.Error("Expected a non-empty array")
	}

	fmt.Println(res)
}

func TestAnkifyUrl(t *testing.T) {

	// Arrange
	const url string = "http://www.paulgraham.com/read.html"
	var anki_text_map = make(map[int]string)
	anki_text_map[1] = url

	godotenv.Load("../../.env")

	// Act
	url_body, _ := docparser.ParseUrl(url)
	res, err := Ankify(url_body, 10)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if res.Questions == nil {
		t.Error("Expected a non-empty array")
	}

	fmt.Println(res)
}
