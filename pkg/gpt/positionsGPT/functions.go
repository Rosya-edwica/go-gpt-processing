package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetFunctionsForPosition(name string) (functions []string, err error) {
	question := fmt.Sprintf(`Составь список из 10 профессиональных функций профессии "%s". Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	resp := gpt.SendRequestToGPT(question)
	functions = strings.Split(resp.Answer, ",")
	if len(functions) <= 1 {
		return nil, gpt.WrongAnswerError
	}
	return
}
