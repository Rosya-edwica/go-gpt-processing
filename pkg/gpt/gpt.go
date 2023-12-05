package gpt

import (
	"context"
	"errors"
	"fmt"
	"go-gpt-processing/pkg/logger"
	"os"
	"regexp"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const (
	QuestionTokenPrice = 0.03 / 1000
	AnswerTokenPrice   = 0.06 / 1000
)

var WrongAnswerError = errors.New("Wrong answer")

type GptResponse struct {
	Question       string
	Answer         string
	ExecutionTime  int64
	QuestionTokens int
	AnswerTokens   int
	TotalTokens    int
	Cost           float64
	Error          error
}

func SendRequestToGPT(query string) GptResponse {
	startTime := time.Now().Unix()

	client := openai.NewClient(os.Getenv("GPT_TOKEN"))
	gptContext, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	completion, err := client.CreateChatCompletion(
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
	if err != nil {
		return GptResponse{Error: err}
	}

	exTime := time.Now().Unix() - startTime
	answer := completion.Choices[0].Message.Content

	if strings.Contains(strings.ToLower(answer), "извините") || len(answer) == 0 {
		return GptResponse{Error: errors.New("GPT не знает что ответить")}
	}
	response := GptResponse{
		Question:       query,
		Answer:         answer,
		ExecutionTime:  exTime,
		QuestionTokens: completion.Usage.PromptTokens,
		AnswerTokens:   completion.Usage.CompletionTokens,
		TotalTokens:    completion.Usage.TotalTokens,
		Cost:           calculateCost(completion.Usage),
	}
	logger.LogInfo.Printf("gpt:\tВремя: %d сек.\tЦена: %f$\tТокены: %d\tТокены(вопрос): %d\tТокены(ответ): %d\tВопрос: %s\t...Ответ: %s\t...",
		response.ExecutionTime, response.Cost, response.TotalTokens, response.QuestionTokens, response.AnswerTokens, response.Question, response.Answer)
	return response
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

func calculateCost(usage openai.Usage) (cost float64) {
	questionTokens := usage.PromptTokens
	answerTokens := usage.CompletionTokens

	cost = float64(questionTokens)*QuestionTokenPrice + float64(answerTokens)*AnswerTokenPrice
	return
}
