package skillsGPT

import (
	"fmt"
	"go-gpt-processing/internal/models"
	"go-gpt-processing/pkg/gpt"
)

func CheckSkillsForDuplicates(skillsPair models.Skill) (isDuplicate bool, err error) {
	question := fmt.Sprintf("Можно ли считать эти навыки дубликатами: '%s' и '%s'? Ответь Да или Нет.", skillsPair.Name, skillsPair.DuplicateName)
	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return false, err
	}
	isDuplicate, err = gpt.ConvertAnswerToBoolean(resp.Answer)
	return
}
