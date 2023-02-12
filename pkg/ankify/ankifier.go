package ankify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type OpenAIRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type AnkiQuestions struct {
	Questions []AnkiQuestion `json:"questions"`
}

type AnkiQuestion struct {
	Question string
	Answer   string
}

const API_URL = "https://api.openai.com/v1/completions"
const ANKI_TEXT = `Make 5 Anki cards for the following text, give it to me in the following 
				format: Q: A: 

				Divide each card with a new line and a ";" at the end of each question and answer. Add 
				more context to each question by searching the web for the answer.`

func GetAnkiCards(anki_text string) (string, error) {

	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("POST", API_URL, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var API_KEY string = os.Getenv("OPENAI_KEY")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	json_request := OpenAIRequest{
		Model:       "text-davinci-003",
		Prompt:      ANKI_TEXT + anki_text,
		MaxTokens:   2048,
		Temperature: 0.5,
	}

	// Attach the JSON payload to the request body
	json_data, err := json.Marshal(json_request)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(json_data))

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Unmarshal the response body into a struct
	var response OpenAIResponse
	json.Unmarshal(body, &response)

	// Check if the response is valid
	if len(response.Choices) == 0 {
		fmt.Println("No response from OpenAI")
		return "", err
	}
	return response.Choices[0].Text, nil
}

// ParseAnkiText parses an Anki text and returns a struct with the questions and answers
func ParseAnkiText(anki_text string) (AnkiQuestions, error) {
	/**
		anki_text: string
		returns: string

		Example:
		anki_text = "Q: What is the capital of France? A: Paris; Q: What is the capital of Germany? A: Berlin"
		anki_cards = ParseAnkiText(anki_text)
		anki_cards.Questions[0].Question = "What is the capital of France?"
		anki_cards.Questions[0].Answer = "Paris"
	**/
	// Split the text based on "line breaks"
	anki_cards := AnkiQuestions{}

	// Split the text based on "line breaks"
	anki_cards_split := strings.Split(anki_text, ";")

	// Loop through the split text and create a new AnkiQuestion for each card
	for _, card := range anki_cards_split {
		// Split the card based on "Q: " and "A: "
		card_split := strings.Split(card, "A: ")

		card_split[0] = strings.TrimPrefix(card_split[0], "Q: ")
		card_split[0] = strings.TrimSpace(card_split[0])
		card_split[1] = strings.TrimSpace(card_split[1])

		// Create a new AnkiQuestion
		anki_card := AnkiQuestion{
			Question: card_split[0],
			Answer:   card_split[1],
		}
		// Append the new AnkiQuestion to the AnkiQuestions struct
		anki_cards.Questions = append(anki_cards.Questions, anki_card)
	}
	return anki_cards, nil
}

func Ankify(anki_text map[int]string) (string, error) {

	anki_questions := ""
	for _, text := range anki_text {

		anki_response, err := GetAnkiCards(text)

		if err != nil {
			fmt.Println(err)
			return "", err
		}
		anki_questions += anki_response
	}

	return anki_questions, nil
}
