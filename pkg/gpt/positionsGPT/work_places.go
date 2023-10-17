package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetWorkPlacesForPosition(name string) (workPlaces []string, timeEx int64, err error) {
	question := fmt.Sprintf(`Составь список из 6 мест где может работать %s . Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	answer, timeEx, err := gpt.SendRequestToGPT(question)
	if err != nil {
		return nil, 0, err
	}

	workPlaces = strings.Split(answer, ",")
	if len(workPlaces) <= 1 || answer == "" || strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, 0, WrongAnswerError
	}
	return
}
