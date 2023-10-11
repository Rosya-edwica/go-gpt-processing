package skillsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"go-gpt-processing/pkg/logger"
)

func GetDescriptionForSkill(name string) (description string, err error) {
	question := fmt.Sprintf("Опиши этот навык - '%s'", name)
	description, err = gpt.SendRequestToGPT(question)
	if err != nil {
		for {
			description, err = (gpt.SendRequestToGPT(question))
			if err == nil {
				break
			}
		}
	}

	logger.Log.Printf("Ответ '%s' для вопроса: %s", description, question)
	return
}
