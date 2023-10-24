package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetOtherNamesForPosition(name string) (otherNames []string, err error) {
	question := fmt.Sprintf(`Составь список из 20 вариантов написания профессии "%s". Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	resp := gpt.SendRequestToGPT(question)
	otherNames = strings.Split(resp.Answer, ",")
	if len(otherNames) <= 1 {
		return nil, gpt.WrongAnswerError
	}
	return
}
