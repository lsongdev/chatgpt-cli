package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/song940/openai-go/openai"
)

var (
	sysrole string
	prompt  string
	rawMode bool
)

func isTermColorSupported() bool {
	return os.Getenv("TERM") == "xterm-color" || os.Getenv("TERM") == "xterm-256color" || os.Getenv("COLORTERM") == "truecolor"
}

func main() {
	defaultRolePrompt := "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible."
	flag.StringVar(&sysrole, "s", defaultRolePrompt, "system role")
	flag.StringVar(&prompt, "p", "", "prompt to use")
	flag.BoolVar(&rawMode, "r", false, "disable color output")
	flag.Parse()
	api := os.Getenv("OPENAI_API")
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("[!] Missing OPENAI_API_KEY, you can find or create your API key here: https://platform.openai.com/account/api-keys")
		return
	}
	if api == "" {
		api = "https://api.openai.com/v1"
	}
	config := openai.Configuration{
		API:    api,
		APIKey: apiKey,
	}
	if !rawMode {
		sysrole = sysrole + "\n(NOTE: Always output as markdown)"
	}
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: sysrole,
		},
	}

	client, err := openai.NewClient(config)
	if err != nil {
		log.Fatal(err)
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

		if isTermColorSupported() && !rawMode {
			rendered, _ := glamour.RenderWithEnvironmentConfig(message)
			return rendered
		} else {
			return message
		}
	}

	if prompt != "" {
		fmt.Println(ask(prompt))
		return
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
