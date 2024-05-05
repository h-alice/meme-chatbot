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

func main() {
	// # Main function
	//
	// This function is the entry point of the program.
	//
	// It reads the user input and sends it to the model.
	//
	// The model response is then printed to the console.
	//
	// The program runs until the user types "exit".
	for {

		var prompt string
		fmt.Print("You: ")
		fmt.Scanln(&prompt)

		if prompt == "exit" {
			break
		}

		// Create the generation parameters
		llmGenerationParameters := LlmGenerationParameters{
			ModelName:     "gpt2",
			TopK:          40,
			TopP:          0.95,
			RepeatPenalty: 1.1,
			Temperature:   0.8,
			Stream:        false,
			MaxTokens:     16,
		}

		// Set the prompt
		llmGenerationParameters.SetPrompt(FormatPrompt(prompt))

		// TODO: Send the prompt to the model
		log.Print(llmGenerationParameters.ToJSON())
	}
}
