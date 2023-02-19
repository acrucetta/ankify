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

func SplitTextIntoRequests(text string, max_tokens int) []string {
	var requests []string
	var num_of_splits int = 1
	var input_token_size int = GetTokenSize(text)

	if input_token_size > max_tokens {
		num_of_splits = (input_token_size / max_tokens) + 1
		for i := 0; i < num_of_splits; i++ {
			startIndex := i * max_tokens
			endIndex := (i + 1) * max_tokens
			requests = append(requests, text[startIndex:endIndex])
		}
	} else {
		requests = append(requests, text)
	}
	log.Printf(`Splitting request into %d parts, the max number of words per request is %d. The document has %d words.`, num_of_splits, max_tokens, input_token_size)
	return requests
}

func SummarizeRequests(requests []string, summarySize int) (string, error) {
	var requestsSummaries string = ""
	if len(requests) > 1 {
		log.Println("Each summary will be approximately", summarySize, "words.")
		for i, request := range requests {
			// Create the prompt for OpenAI
			summaryPrompt := strings.Replace(SUMMARY_HELPER, "{text}", request, 1)
			summaryPrompt = strings.Replace(summaryPrompt, "{summary_size}", strconv.Itoa(summarySize), 1)
			ankiResponse, err := CallOpenAI(summaryPrompt)
			if err != nil {
				return "", err
			}
			// Log the request number using logger
			log.Printf("Finished processing request %d of %d.", i+1, len(requests))
			requestsSummaries += ankiResponse
		}
	} else {
		requestsSummaries = requests[0]
	}
	return requestsSummaries, nil
}

func CreateAnkiCards(text string, cardNum int) (AnkiQuestions, error) {
	ankiQuestions := AnkiQuestions{}
	anki_token_size := GetTokenSize(text)
	log.Printf("The summary has %d tokens.", anki_token_size)
	if anki_token_size > 2048 {
		log.Fatal("The summary is too long, we will use only the first 2000 words.")
		text = text[:2000]
	}
	ankiPrompt := strings.Replace(QUESTION_HELPER, "{text}", text, 1)
	ankiPrompt = strings.Replace(ankiPrompt, "{card_num}", strconv.Itoa(cardNum), 1)
	ankiResponse, err := CallOpenAI(ankiPrompt)
	if err != nil {
		return AnkiQuestions{}, err
	}
	parsedQuestions, err := ParseAnkiText(ankiResponse)
	if err != nil {
		return AnkiQuestions{}, err
	}
	log.Println("Successfully parsed the anki cards, adding them to the CSV.")
	ankiQuestions.Questions = append(ankiQuestions.Questions, parsedQuestions.Questions...)
	return ankiQuestions, nil
}

func Ankify(ankiText map[int]string, cardNum int) (AnkiQuestions, error) {
	ankiQuestions := AnkiQuestions{}
	const max_tokens int = 2048
	for _, text := range ankiText {
		// Check the number of tokens in the text
		// doesn't exceed the maximum number of tokens
		// allowed by OpenAI (3800); if it does, split
		// the text into multiple requests
		requests := SplitTextIntoRequests(text, max_tokens)
		var summary_size int = max_tokens / len(requests)
		summarized_text, err := SummarizeRequests(requests, summary_size)
		if err != nil {
			log.Fatal(err)
			return AnkiQuestions{}, err
		}

		// Create the anki cards from the summary
		ankiQuestionsForText, err := CreateAnkiCards(summarized_text, cardNum)
		if err != nil {
			log.Fatal(err)
			return AnkiQuestions{}, err
		}
		ankiQuestions.Questions = append(ankiQuestions.Questions, ankiQuestionsForText.Questions...)
	}
	return ankiQuestions, nil
}

func GetTokenSize(text string) int {
	// We assume a token is about 4 characters
	// and we count the number of spaces
	// to get the number of tokens
	return len(text) / 4
}
