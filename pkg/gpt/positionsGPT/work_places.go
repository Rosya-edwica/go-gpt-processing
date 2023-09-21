package positionsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetWorkPlacesForPosition(name string) (workPlaces []string, err error) {
	question := fmt.Sprintf(`Составь список из 6 мест где может работать %s . Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	answer, err := gpt.SendRequestToGPT(question)
	workPlaces = strings.Split(answer, ",")

	if len(workPlaces) <= 1 {
		return nil, errors.New(fmt.Sprintf("Не удалось поделить ответ по запятым: %s", answer))
	}
	if answer == "" {
		return nil, errors.New(fmt.Sprintf("Пустое ответ для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", answer, name))
	}
	return
}
