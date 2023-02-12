package ankify

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestAnkify(t *testing.T) {

	// Arrange
	const anki_text string = "This is a test"
	var anki_text_map = make(map[int]string)
	anki_text_map[1] = anki_text

	godotenv.Load("../../.env")

	// Act
	res, err := Ankify(anki_text_map)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if res == "" {
		t.Error("Expected a non-empty string")
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
