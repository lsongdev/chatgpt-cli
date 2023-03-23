package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/song940/openai-go/openai"
)

func main() {

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("[!] Missing OPENAI_API_KEY, you can find or create your API key here: https://platform.openai.com/account/api-keys")
		return
	}
	config := openai.Configuration{
		APIKey: apiKey,
	}
	client, err := openai.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	var prompt = "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible."
	if len(os.Args) > 1 {
		prompt = os.Args[1]
	}
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: prompt,
		},
	}
	ask := func(text string) string {
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
		res, err := client.CreateChatCompletion(request)
		if err != nil {
			return fmt.Sprintf("%v", err)
		}
		message := res.Choices[0].Message.Content
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    "assistant",
				Content: message,
			},
		)
		rendered, _ := glamour.RenderWithEnvironmentConfig(message)
		return rendered
	}
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		input, _ := io.ReadAll(os.Stdin)
		fmt.Println(ask(string(input)))
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You> ")
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		out := ask(text)
		fmt.Println("OpenAI>")
		fmt.Println(out)
	}

}
