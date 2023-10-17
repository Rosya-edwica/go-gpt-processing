package skillsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
)

func GetDescriptionForSkill(name string) (description string, timeEx int64, err error) {
	question := fmt.Sprintf("Опиши этот навык - '%s'", name)
	description, timeEx, err = gpt.SendRequestToGPT(question)
	return
}
