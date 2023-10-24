package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetWorkPlacesForPosition(name string) (workPlaces []string, err error) {
	question := fmt.Sprintf(`Составь список из 6 мест где может работать %s . Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	resp := gpt.SendRequestToGPT(question)
	workPlaces = strings.Split(resp.Answer, ",")
	if len(workPlaces) <= 1 {
		return nil, gpt.WrongAnswerError
	}
	return
}
