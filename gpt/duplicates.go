package gpt

import (
	"fmt"
	"gpt-skills/logger"
	"gpt-skills/models"
)


func CheckSkillsForDuplicates(skillsPair *models.Pair) (err error) {
	var answer string
	question := fmt.Sprintf("Можно ли считать эти навыки дубликатами: '%s' и '%s'? Ответь Да или Нет.", skillsPair.First, skillsPair.Second)
	answer, err =  SendRequestToGPT(question)
	if err != nil {
		for {
			fmt.Println(err)
			answer, err = (SendRequestToGPT(question))
			if err == nil {
				break
			}
		}
	}
	skillsPair.IsDuplicate, err = convertAnswerToBoolean(answer)
	logger.Log.Printf("Ответ '%s' для вопроса: %s", answer, question)
	return
}

