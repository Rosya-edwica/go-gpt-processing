package positionsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetDescriptionForPosition(name string) (descr string, err error) {
	question := fmt.Sprintf(`Выступи в роли карьерного эксперта и сделай описание профессии "%s". Описание должно быть написано простым языком, для людей, которые могут не знать об этой профессии. Напиши подробное описание не больше 2 абзацев. Количество символов не больше 1000`, name)
	descr, err = gpt.SendRequestToGPT(question)
	if descr == "" {
		return "", errors.New(fmt.Sprintf("Пустое описание для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(descr), "я не могу") {
		return "", errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", descr, name))
	}
	return
}
