package positionsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetOtherNamesForPosition(name string) (otherNames []string, timeEx int64, err error) {
	question := fmt.Sprintf(`Составь список из 20 вариантов написания профессии "%s". Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	answer, timeEx, err := gpt.SendRequestToGPT(question)
	otherNames = strings.Split(answer, ",")

	if len(otherNames) <= 1 {
		return nil, timeEx, errors.New(fmt.Sprintf("Не удалось поделить ответ по запятым: %s", answer))
	}
	if answer == "" {
		return nil, timeEx, errors.New(fmt.Sprintf("Пустое описание для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, timeEx, errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", answer, name))
	}
	return
}
