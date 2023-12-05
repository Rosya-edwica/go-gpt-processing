package positionsGPT

import (
	"fmt"
	"go-gpt-processing/internal/gpt"
	"go-gpt-processing/internal/models"
	"strings"
)

func GetDisabilityForPosition(name string, disabilities []models.Disability) (disability []models.Disability, err error) {
	question := fmt.Sprintf(`Я занимаюсь подбором профессий для людей с ограниченными возможностями. 
	Ниже будет представлен список инвалидностей. Тебе нужно выбрать из списка только те, 
	которые не являются препятствием для освоения профессии "%s". 
	Если для данной профессии ни одна группа инвалидности не подходит, напиши "нельзя"
	1. Ограничения по зрению 
	 2. Ограничения по слуху 
	 3. Ограничения опорно-двигательного аппарата верхних конечностей 
	 4. Ограничения опорно-двигательного аппарата нижних конечностей
	 5. Нарушения речи 
	 6. Нарушения поведения и общения, интеллектуальных процессов`, name)
	response := gpt.SendRequestToGPT(question)
	if response.Error != nil {
		return nil, response.Error
	}
	if strings.Contains(strings.ToLower(response.Answer), "нельзя") {
		fmt.Printf("Нельзя для: '%s' -> %s \n", name, response.Answer)
		return nil, nil
	}

	return convertAnswerToDisabilities(response.Answer, disabilities)
}

func convertAnswerToDisabilities(text string, disabilities []models.Disability) (positionDisabilities []models.Disability, err error) {
	items := reLines.FindAllString(text, -1)
	if len(items) <= 1 {
		return nil, gpt.WrongAnswerError
	}

	for _, i := range items {
		item := rePointDigit.ReplaceAllString(i, "")
		item = strings.ReplaceAll(item, ".", "")
		item = strings.TrimSpace(item)
		for _, dis := range disabilities {
			if item == dis.Name {
				positionDisabilities = append(positionDisabilities, dis)
			}
		}
	}

	return positionDisabilities, nil
}
