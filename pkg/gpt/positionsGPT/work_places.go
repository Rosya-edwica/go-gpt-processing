package positionsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetWorkPlacesForPosition(name string) (workPlaces []string, timeEx int64, err error) {
	question := fmt.Sprintf(`Составь список из 6 мест где может работать %s . Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	answer, timeEx, err := gpt.SendRequestToGPT(question)
	workPlaces = strings.Split(answer, ",")

	if len(workPlaces) <= 1 {
		return nil, timeEx, errors.New(fmt.Sprintf("Не удалось поделить ответ по запятым: %s", answer))
	}
	if answer == "" {
		return nil, timeEx, errors.New(fmt.Sprintf("Пустое ответ для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, timeEx, errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", answer, name))
	}
	return
}
