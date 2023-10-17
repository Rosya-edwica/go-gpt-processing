package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetEducationForPosition(name string) (education []string, timeEx int64, err error) {
	question := fmt.Sprintf(`
		Каким образованием необходимо владеть специалисту "%s" России. Выбери из списка:
		1. Без образования
		2. Среднее профессиональное образование
		3. Высшее образование
		Если данному специалисту необходимо только высшее образование, то поставь "Высшее образование". 
		Если данного специалиста обучают и в заведениях высшего образования, и в заведениях среднего профессионального образования, 
		то поставь два варианта "Среднее профессиональное образование" и "Высшее образование". Если данной профессии можно научиться и без образования, 
		то поставь "Без образования". Если данной профессии можно обучиться где угодно и даже самостоятельно, то поставь все три ответа "Без образования", 
		"Среднее профессиональное образование" и "Высшее образование". Запиши ответы без ковычек, в одну строку, через знак запятая. 
		Учитывай только специфику Российского образования.
	`, name)
	answer, timeEx, err := gpt.SendRequestToGPT(question)
	if !strings.Contains("без образования", strings.ToLower(answer)) &&
		!strings.Contains("среднее профессиональное образование", strings.ToLower(answer)) &&
		!strings.Contains("высшее образование", strings.ToLower(answer)) {
		return []string{}, timeEx, WrongAnswerError
	}
	education = strings.Split(answer, ",")
	if err != nil {
		return []string{}, timeEx, WrongAnswerError
	}
	return
}
