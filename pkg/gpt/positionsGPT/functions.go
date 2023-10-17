package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetFunctionsForPosition(name string) (functions []string, timeEx int64, err error) {
	question := fmt.Sprintf(`Составь список из 10 профессиональных функций профессии "%s". Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	answer, timeEx, err := gpt.SendRequestToGPT(question)
	if err != nil {
		return nil, 0, err
	}

	functions = strings.Split(answer, ",")
	if len(functions) <= 1 || answer == "" || strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, 0, WrongAnswerError
	}
	return
}
