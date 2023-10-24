package skillsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
)

func GetDescriptionForSkill(name string) (description string, err error) {
	question := fmt.Sprintf("Опиши этот навык - '%s'", name)
	resp := gpt.SendRequestToGPT(question)
	return resp.Answer, resp.Error
}
