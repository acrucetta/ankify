package ankify

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"log"
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
	Tag      string
}

const API_URL = "https://api.openai.com/v1/completions"
const QUESTION_HELPER = `Make {card_num} Anki cards for the following text, give it to me in the following format: Q: [Insert question here] A: [Insert answer here] \n\n The text is the following, be detailed and include the most unique and helpful points: {text}`
const SUMMARY_HELPER = `Summarize the following text into less than {summary_size} words in detail with the most unique and helpful points, explain it to me like an expert: {text}`

func CallOpenAI(prompt string) (string, error) {

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("POST", API_URL, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	var API_KEY string = os.Getenv("OPENAI_KEY")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	json_request := OpenAIRequest{
		Model:       "text-davinci-003",
		Prompt:      prompt,
		MaxTokens:   2048,
		Temperature: 0.5,
	}

	// Attach the JSON payload to the request body
	json_data, err := json.Marshal(json_request)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(json_data))

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Unmarshal the response body into a struct
	var response OpenAIResponse
	json.Unmarshal(body, &response)

	// Check if the response is valid
	if len(response.Choices) == 0 {
		log.Fatal("No response from OpenAI")
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

	// Split the text based on "\n\n")
	anki_cards_split := strings.Split(anki_text, "\n\n")

	// Loop through the split text and create a new AnkiQuestion for each card
	for _, card := range anki_cards_split {

		if card == "" {
			continue
		}

		// If the card doesn't contain "Q: " or "A: ", skip it
		if !strings.Contains(card, "Q: ") || !strings.Contains(card, "A: ") {
			continue
		}

		// Split the card based on "Q: " and "A: "
		card_split := strings.Split(card, "A:")

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

func Ankify(anki_text map[int]string, card_num int) (AnkiQuestions, error) {

	anki_questions := AnkiQuestions{}
	for _, text := range anki_text {

		// Check the number of tokens in the text
		// doesn't exceed the maximum number of tokens
		// allowed by OpenAI (3800); if it does, split
		// the text into multiple requests
		var requests []string
		var num_splits int = 1
		var token_size int = GetTokenSize(text)
		const MAX_TOKENS = 2048

		if len(text) > MAX_TOKENS {
			num_splits = token_size / MAX_TOKENS
			for i := 0; i < num_splits; i++ {
				// Get the start and end index of the text
				start_index := i * MAX_TOKENS
				end_index := (i + 1) * MAX_TOKENS
				requests = append(requests, text[start_index:end_index])
			}
		} else {
			requests = append(requests, text)
		}
		log.Printf(`Splitting request into %d parts, the max number of words per request is 2048. The document has %d words.`, num_splits, token_size)

		// Loop through the requests and summarize each request
		// using OpenAI
		var requests_summaries string = ""

		if len(requests) > 1 {
			var summary_size int = 2000 / num_splits
			log.Println("Each summary will be approximately", summary_size, "words.")
			for i, request := range requests {

				// Create the prompt for OpenAI
				summary_prompt := strings.Replace(SUMMARY_HELPER, "{text}", request, 1)
				summary_prompt = strings.Replace(summary_prompt, "{summary_size}", strconv.Itoa(summary_size), 1)
				anki_response, err := CallOpenAI(summary_prompt)
				if err != nil {
					log.Fatal(err)
					return AnkiQuestions{}, err
				}
				// Log the request number using logger
				log.Printf("Finished processing request %d of %d.", i+1, len(requests))
				requests_summaries += anki_response
			}
		} else {
			requests_summaries = requests[0]
		}

		// Call the OpenAI API to create the anki cards from the summary
		anki_token_size := GetTokenSize(requests_summaries)
		log.Printf("The summary has %d words.", anki_token_size)
		anki_prompt := strings.Replace(QUESTION_HELPER, "{text}", requests_summaries, 1)
		anki_prompt = strings.Replace(anki_prompt, "{card_num}", strconv.Itoa(card_num), 1)
		anki_response, err := CallOpenAI(anki_prompt)
		if err != nil {
			log.Fatal(err)
			return AnkiQuestions{}, err
		}
		parsed_questions, err := ParseAnkiText(anki_response)
		if err != nil {
			log.Fatal(err)
			return AnkiQuestions{}, err
		}
		anki_questions.Questions = append(anki_questions.Questions, parsed_questions.Questions...)
	}
	return anki_questions, nil
}

func GetTokenSize(text string) int {
	return len(strings.Split(text, " "))
}
