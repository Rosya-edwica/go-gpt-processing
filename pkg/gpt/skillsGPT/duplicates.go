package skillsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"go-gpt-processing/pkg/models"
)

func CheckSkillsForDuplicates(skillsPair models.Skill) (isDuplicate bool, exTime int64, err error) {
	var answer string
	question := fmt.Sprintf("Можно ли считать эти навыки дубликатами: '%s' и '%s'? Ответь Да или Нет.", skillsPair.Name, skillsPair.DuplicateName)
	answer, exTime, err = gpt.SendRequestToGPT(question)
	if err != nil {
		return false, 0, err
	}
	isDuplicate, err = gpt.ConvertAnswerToBoolean(answer)
	return
}
