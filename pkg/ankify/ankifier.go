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

type ResponseBody struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Message      ResponseMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
	Index        int             `json:"index"`
}

type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string           `json:"model"`
	Messages []RequestMessage `json:"messages"`
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnkiQuestions struct {
	Questions []AnkiQuestion `json:"questions"`
}

type AnkiQuestion struct {
	Question string
	Answer   string
	Tag      string
}

const API_URL = "https://api.openai.com/v1/chat/completions"
const QUESTION_HELPER = `Assume you’re an expert learning and memory. You will help me write Anki questions. The Anki questions should be concise and focus on one atomic unit. They should encode ideas from multiple angles, connect and relate ideas, and unambiguously produce a specific answer. Additionally, the questions must make clear what shape of answer is expected and ensure reviewers retrieve answers from memory. Avoid yes-no questions.

I want you to make {card_num} Anki cards for the following text, give it to me in the following format: 

Q: [Insert question here] 
A: [Insert answer here] 
\n\n 

The text is the following: 
{text}`
const SUMMARY_HELPER = `Assume you’re an expert in summarizing text to the most important points of paragraph in a way that retains the original meaning and context of the pragraph, I want you to summarize the following text into less than {summary_size} words with the most unique and helpful points: {text}`

func CallOpenAI(prompt string) (string, error) {

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("POST", API_URL, nil)
	if err != nil {
		log.Fatal(err, os.Stderr)
		return "", err
	}

	var API_KEY string = os.Getenv("OPENAI_KEY")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	// Create a message struct
	message := RequestMessage{
		Role:    "user",
		Content: prompt,
	}

	json_request := OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: []RequestMessage{message},
	}

	// Attach the JSON payload to the request body
	json_data, err := json.Marshal(json_request)

	if err != nil {
		log.Fatal(err, os.Stderr)
		return "", err
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(json_data))
	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, os.Stderr)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err, os.Stderr)
		return "", err
	}

	// Unmarshal the response body into a struct
	var response ResponseBody
	json.Unmarshal(body, &response)

	// Check if the response is valid
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	} else {
		log.Fatal("The error response is: ", string(body))
		return "", nil
	}
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

func CreateAnkiCards(text string, card_num int) (AnkiQuestions, error) {
	anki_questions := AnkiQuestions{}
	anki_token_size := GetTokenSize(text)
	log.Printf("The summary has %d tokens.", anki_token_size)
	if anki_token_size > 3000 {
		log.Printf("The summary is too long, we will use only the first 3000 tokens.")
		text = text[:(3000 * 4)]
	}
	anki_prompt := strings.Replace(QUESTION_HELPER, "{text}", text, 1)
	anki_prompt = strings.Replace(anki_prompt, "{card_num}", strconv.Itoa(card_num), 1)
	log.Printf("The final length of the prompt is %d tokens.", GetTokenSize(anki_prompt))
	anki_response, err := CallOpenAI(anki_prompt)
	if err != nil {
		return AnkiQuestions{}, err
	}
	parsed_questions, err := ParseAnkiText(anki_response)
	if err != nil {
		return AnkiQuestions{}, err
	}
	log.Println("Successfully parsed the anki cards, adding them to the CSV.")
	anki_questions.Questions = append(anki_questions.Questions, parsed_questions.Questions...)
	return anki_questions, nil
}

func Ankify(ankiText map[int]string, cardNum int) (AnkiQuestions, error) {
	ankiQuestions := AnkiQuestions{}
	const max_tokens int = 3000
	for _, text := range ankiText {
		// Check the number of tokens in the text
		// doesn't exceed the maximum number of tokens
		// allowed by OpenAI (3800); if it does, split
		// the text into multiple requests
		requests := SplitTextIntoRequests(text, max_tokens)
		var summary_size int = max_tokens / len(requests)
		summarized_text, err := SummarizeRequests(requests, summary_size)
		if err != nil {
			log.Fatal(err, os.Stderr)
			return AnkiQuestions{}, err
		}

		// Create the anki cards from the summary
		ankiQuestionsForText, err := CreateAnkiCards(summarized_text, cardNum)
		if err != nil {
			log.Fatal(err, os.Stderr)
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
