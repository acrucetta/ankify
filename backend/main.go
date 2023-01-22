package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const API_URL = "https://api.openai.com/v1/completions"
const ANKI_TEXT = "Make an Anki card for the following text, give it to me in the following format Question: Answer:"

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

func openAIHandler(c *gin.Context) {
	var request map[string]string
	c.BindJSON(&request)
	prompt := request["prompt"]

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("POST", API_URL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var API_KEY string = os.Getenv("OPENAI_KEY")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	json_request := OpenAIRequest{
		Model:       "text-davinci-003",
		Prompt:      ANKI_TEXT + prompt,
		MaxTokens:   2048,
		Temperature: 0.5,
	}

	// Attach the JSON payload to the request body
	json_data, err := json.Marshal(json_request)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(json_data))

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Unmarshal the response body into a struct
	var response OpenAIResponse
	json.Unmarshal(body, &response)

	// Check if the response is valid
	if len(response.Choices) == 0 {
		fmt.Println("No response from OpenAI")
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": response.Choices[0].Text})
}

func main() {
	router := gin.Default()
	router.POST("/getquestions", openAIHandler)

	// Allow CORS Headers
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type"}}))
	router.Run(":8080")
}
