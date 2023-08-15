package gpt

import (
	"errors"
	"fmt"
	"gpt-skills/logger"
	"gpt-skills/models"
	"strings"
	"time"
)


func CheckSkillsForTypeGroup(skill *models.Skill) (err error) {
	question := fmt.Sprintf(`Если слово «%s» является знанием, умением, квалификацией, навыком, способностью, профнавыком, практическим навыком, компетенцией, то поставь – 1, если это слово является профессией, должностью, работой, квалификацией, то поставь – 2, в ином случае поставь - 3`, skill.Name)
	answer, err := SendRequestToGPT(question)
	if err != nil {
		fmt.Println("ОШИБКА:", err)
		time.Sleep(5 * time.Second)
		return
	}
	fmt.Println(answer)
	if strings.Contains(answer, "профессия/специальность/должность") || strings.Contains(answer, "2") {
		skill.Group = "профессия"
	} else if strings.Contains(answer, "навык") || strings.Contains(answer, "1") {
		skill.Group = "навык"
	} else if strings.Contains(answer, "другое") || strings.Contains(answer, "3") {
		skill.Group = "другое"
	} else {
		return errors.New(fmt.Sprintf("ОШИБКА: ответ - %s. вопрос: %s", answer, question))
	}
	
	logger.Log.Printf("Ответ '%s' для вопроса: %s", answer, question)

	return
}