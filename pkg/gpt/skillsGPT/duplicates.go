package skillsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"go-gpt-processing/pkg/logger"
	"go-gpt-processing/pkg/models"
)

func CheckSkillsForDuplicates(skillsPair *models.Skill) (err error) {
	var answer string
	question := fmt.Sprintf("Можно ли считать эти навыки дубликатами: '%s' и '%s'? Ответь Да или Нет.", skillsPair.Name, skillsPair.DuplicateName)
	answer, err = gpt.SendRequestToGPT(question)
	if err != nil {
		for {
			fmt.Println(err)
			answer, err = (gpt.SendRequestToGPT(question))
			if err == nil {
				break
			}
		}
	}
	skillsPair.IsDuplicate, err = gpt.ConvertAnswerToBoolean(answer)
	logger.Log.Printf("Ответ '%s' для вопроса: %s", answer, question)
	return
}
