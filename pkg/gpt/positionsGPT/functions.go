package positionsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetFunctionsForPosition(name string) (functions []string, timeEx int64, err error) {
	question := fmt.Sprintf(`Составь список из 10 профессиональных функций профессии "%s". Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	answer, timeEx, err := gpt.SendRequestToGPT(question)
	functions = strings.Split(answer, ",")
	if err != nil {
		fmt.Println("Ошибка при подборе функций для вопроса: ", question)
		return []string{}, timeEx, err
	}
	if len(functions) <= 1 {
		return nil, timeEx, errors.New(fmt.Sprintf("Не удалось поделить ответ по запятым: %s", answer))
	}
	if answer == "" {
		return nil, timeEx, errors.New(fmt.Sprintf("Нет функций для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, timeEx, errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", answer, name))
	}
	return
}
