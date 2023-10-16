package positionsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetAboutForPosition(name string) (about string, timeEx int64, err error) {
	question := fmt.Sprintf(`Составь описание профессии "%s" в одну строчку`, name)
	about, timeEx, err = gpt.SendRequestToGPT(question)
	if about == "" {
		return "", timeEx, errors.New(fmt.Sprintf("Пустое описание для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(about), "я не могу") {
		return "", timeEx, errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", about, name))
	}
	return
}
