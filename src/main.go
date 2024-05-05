package main

import (
	"fmt"
	"log"
)

const CHAT_TEMPLATE = `<start_of_turn>user
%s<end_of_turn>
<start_of_turn>model
`
const CHAT_TEMPLATE_END = "<end_of_turn>"

type LlmGenerationParameters struct {
	MaxTokens int `json:"max_tokens"`
}

// # Prompt formatter
//
// This function formats the prompt to be sent to the model.
//
// Parameters:
//
// - prompt: the user prompt
func formatPrompt(prompt string) string {
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

		// Format the prompt
		formattedPrompt := formatPrompt(prompt)

		// TODO: Send the prompt to the model
		log.Print(formattedPrompt)
	}
}
