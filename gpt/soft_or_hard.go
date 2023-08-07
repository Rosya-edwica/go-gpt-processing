package gpt

import (
	"fmt"
	"gpt-skills/logger"
	"gpt-skills/models"
	"strings"
)

func CheckSkillIsSoftOrHard(softOrHard string, skill *models.Skill) (err error) {
	var answer, question string
	if softOrHard == "soft" {
		question = fmt.Sprintf("Ответь Да или Нет - это софт-скилл: '%s'?", skill.Name)
	} else {
		question = fmt.Sprintf("Ответь Да или Нет - это хард-скилл: '%s'?", skill.Name)
	}

	answer, err =  SendRequestToGPT(question)
	if err != nil {
		for {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				return
			}
			answer, err = (SendRequestToGPT(question))
			if err == nil {
				break
			}
		}
	}
	skill.IsValid, err = convertAnswerToBoolean(answer)	
	logger.Log.Printf("Ответ '%s' для вопроса: %s", answer, question)
	return
}
