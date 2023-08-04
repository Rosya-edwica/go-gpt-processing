package gpt

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)


func SendRequestToGPT(query string) (answer string, err error) {
	client := openai.NewClient(os.Getenv("GPT_TOKEN"))
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: query,
				},
			},
		},
	)
	if err != nil {
		panic(err)
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

func checkErr(err error) { 
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}