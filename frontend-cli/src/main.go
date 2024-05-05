package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const CHAT_TEMPLATE = `<start_of_turn>user
%s<end_of_turn>
<start_of_turn>model
`
const CHAT_TEMPLATE_END = "<end_of_turn>"

type LlmGenerationParameters struct {
	ModelName     string  `json:"model"`
	Prompt        string  `json:"prompt"`
	TopK          int     `json:"top_k"`
	TopP          float64 `json:"top_p"`
	RepeatPenalty float64 `json:"repeat_penalty"`
	Temperature   float64 `json:"temperature"`
	Stream        bool    `json:"stream"`
	MaxTokens     int     `json:"max_tokens"`
}

// # Check and fix generation parameters
//
// This function checks the generation parameters and fixes them if needed.
// If a parameter is missing or invalid, it is set to a default value.
//
// The default values are suggested by the llama-cpp-python library.
// Check the `llama_cpp/server/types.py` source code for more information.
func (lgp *LlmGenerationParameters) CheckAndFix() {
	if lgp.TopK <= 0 {
		lgp.TopK = 40
	}
	if lgp.TopP <= 0 || lgp.TopP > 1.0 {
		lgp.TopP = 0.95
	}
	if lgp.RepeatPenalty <= 0 {
		lgp.RepeatPenalty = 1.1
	}
	if lgp.Temperature <= 0 {
		lgp.Temperature = 0.8
	}
	if lgp.MaxTokens <= 0 {
		lgp.MaxTokens = 16
	}
}

// # Set the prompt
//
// This function sets the prompt to be sent to the model.
//
// Note that the SetPrompt return a copy of the generation parameters with the prompt set.
func (lgp LlmGenerationParameters) SetPrompt(prompt string) LlmGenerationParameters {
	lgp.Prompt = prompt
	return lgp
}

// # To JSON
//
// This function converts the generation parameters to a JSON string.
func (lgp *LlmGenerationParameters) ToJSON() string {
	lgp.CheckAndFix()
	jsonData, _ := json.Marshal(lgp)
	return string(jsonData)
}

// # Prompt formatter
//
// This function formats the prompt to be sent to the model.
//
// Parameters:
//
// - prompt: the user prompt
func FormatPrompt(prompt string) string {
	return fmt.Sprintf(CHAT_TEMPLATE, prompt)
}

const SampleResponse = `{
	"id": "cmpl-555e840b-6921-44e8-9f6f-ab9fcd859624",
	"object": "text_completion",
	"created": 1714891381,
	"model": "gpt2",
	"choices": [
		{
			"text": "好",
			"index": 0,
			"logprobs": null,
			"finish_reason": "stop"
		}
	],
	"usage": {
		"prompt_tokens": 12,
		"completion_tokens": 1,
		"total_tokens": 13
	}
}`

type LlmResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage interface{} `json:"usage"`
}

// # Parse response
//
// This function parses the response from the model.
func ParseResponse(response string) LlmResponse {
	var llmResponse LlmResponse
	json.Unmarshal([]byte(response), &llmResponse)
	return llmResponse
}

// # Connect to server endpoint and send prompt
//
// This function connects to the local server and sends the prompt to the model.
func SendPrompt(server string, port int, endpoint string, param_with_prompt LlmGenerationParameters) (string, error) {

	// Construct the URL
	url := fmt.Sprintf("http://%s:%d/%s", server, port, endpoint)

	// Send the prompt to the model
	resp, err := http.Post(url, "application/json", strings.NewReader(param_with_prompt.ToJSON()))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() // Close the response body

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// # Model I/O handler
//
// This function handles the communication between the user prompt queue and the model response queue.
//
// Note that the `user_prompt_queue` expects the prompt has been given by the user.
func modelIoHandler(ctx context.Context, server string, port int, endpoint string, param_with_prompt_queue <-chan LlmGenerationParameters, model_response_queue chan<- string, wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Get the user prompt
			param_with_prompt := <-param_with_prompt_queue

			for {
				// Send the prompt to the model
				response, err := SendPrompt(server, port, endpoint, param_with_prompt)
				if err != nil {
					log.Println(err)
					time.Sleep(1 * time.Second)
					continue
				}

				// Get the actual response from the model
				model_output := ParseResponse(response).Choices[0].Text

				// Send the model response to the response queue
				model_response_queue <- model_output

				break
			}
		}
	}
}

func main() {
	// Test sending a prompt to the model
	server := "backend"
	port := 8000
	endpoint := "v1/completions"
	param_template := LlmGenerationParameters{
		ModelName:     "",
		TopK:          64,
		TopP:          0.9,
		RepeatPenalty: 1.2,
		Temperature:   0.9,
		Stream:        false,
		MaxTokens:     32,
	}

	ctx := context.Background()

	// Create channels.
	param_with_prompt_queue := make(chan LlmGenerationParameters)
	model_response_queue := make(chan string)

	// Create a wait group.
	wg := new(sync.WaitGroup)

	// Add the model I/O handler to the wait group.
	wg.Add(1)

	// Start the model I/O handler.
	go modelIoHandler(ctx, server, port, endpoint, param_with_prompt_queue, model_response_queue, wg)

	// User cli interaction.
	for {
		var user_input string
		fmt.Print("User: ")
		fmt.Scanln(&user_input)

		// Set the prompt
		param_with_prompt := param_template.SetPrompt(FormatPrompt(user_input))

		// Send the prompt to the model
		param_with_prompt_queue <- param_with_prompt

		// Get the model response
		response := <-model_response_queue

		// Print the model response
		fmt.Println("Model:", response)
	}
}
