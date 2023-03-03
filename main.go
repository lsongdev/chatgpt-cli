package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/song940/openai-go/openai"
)

const maxTokens = 4096

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing OPENAI_API_KEY, you can find or create your API key here: https://platform.openai.com/account/api-keys")
	}
	config := openai.Configuration{
		APIKey: apiKey,
	}
	client, err := openai.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible.",
		},
	}
	totalTokens := 0
	for {
		fmt.Printf("(%d/%d)> ", totalTokens, maxTokens)
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    "user",
				Content: text,
			},
		)
		request := openai.ChatCompletionRequest{
			Model:           openai.GPT3_5_Trubo,
			Messages:        messages,
			MaxTokens:       3000,
			NumberOfChoices: 1,
		}
		resp, err := client.CreateChatCompletion(request)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// totalTokens = resp.Usage.TotalTokens
		message := resp.Choices[0].Message.Content
		rendered, _ := glamour.RenderWithEnvironmentConfig(message)
		fmt.Println(rendered)
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    "assistant",
				Content: message,
			},
		)
	}
}
