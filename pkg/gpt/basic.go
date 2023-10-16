package gpt

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const (
	QuestionTokenPrice = 0.03 / 1000
	AnswerTokenPrice   = 0.06 / 1000
)

func SendRequestToGPT(query string) (answer string, exTime int64, err error) {
	startTime := time.Now().Unix()

	client := openai.NewClient(os.Getenv("GPT_TOKEN"))
	gptContext, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	response, err := client.CreateChatCompletion(
		gptContext,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
		},
	)
	exTime = time.Now().Unix() - startTime
	if err != nil {
		return "", exTime, err
	}
	answer = response.Choices[0].Message.Content
	AddCostToAmount(response.Usage)

	if strings.Contains(strings.ToLower(answer), "извините") {
		return "", exTime, errors.New("GPT не знает что ответить")
	}
	return
}

func ConvertAnswerToBoolean(answer string) (bool, error) {
	re := regexp.MustCompile(`yes,|no,|no|yes`)
	answer = strings.ReplaceAll(strings.ToLower(answer), ".", "")
	answer = re.FindString(answer)
	answer = strings.ReplaceAll(answer, ",", "")

	switch answer {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, errors.New(fmt.Sprintf("Неправильный ответ: %s", answer))
	}
}

func AddCostToAmount(usage openai.Usage) {
	questionTokens := usage.PromptTokens
	answerTokens := usage.CompletionTokens

	price := float64(questionTokens)*QuestionTokenPrice + float64(answerTokens)*AnswerTokenPrice

	amount, err := strconv.ParseFloat(os.Getenv("GPT_4_AMOUNT"), 64)
	if err != nil {
		panic(err)
	}
	amount += price
	os.Setenv("GPT_4_AMOUNT", fmt.Sprintf("%f", amount))
}
