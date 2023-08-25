package gpt

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)


func SendRequestToGPT(query string) (answer string, err error) {
	client := openai.NewClient(os.Getenv("GPT_TOKEN"))
	gptContext, cancel := context.WithTimeout(context.Background(), time.Second * 20)
	defer cancel()
	response, err := client.CreateChatCompletion(
		gptContext,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: query,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	answer = response.Choices[0].Message.Content
	return
}

func convertAnswerToBoolean(answer string) (bool, error) {
	re := regexp.MustCompile(`да,|нет,|нет|не|да`)
	answer = strings.ReplaceAll(strings.ToLower(answer), ".", "")
	answer = re.FindString(answer)
	answer = strings.ReplaceAll(answer, ",", "")

	switch answer{
	case "да": return true, nil
	case "нет": return false, nil
	default: return false, errors.New(fmt.Sprintf("Неправильный ответ: %s", answer))
	}
}