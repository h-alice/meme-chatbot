package main

import (
	"encoding/json"
	"fmt"
	"log"
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
func (lgp *LlmGenerationParameters) SetPrompt(prompt string) {
	lgp.Prompt = prompt
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

/*
Sample response body:
{
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
}
*/

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

// # Connect to local server and send prompt
//
// This function connects to the local server and sends the prompt to the model.

func main() {
	// Test parse
	llmResponse := ParseResponse(SampleResponse)
	log.Println(llmResponse)
}
