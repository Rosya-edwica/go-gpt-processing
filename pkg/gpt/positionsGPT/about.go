package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
)

func GetAboutForPosition(name string) (about string, err error) {
	question := fmt.Sprintf(`Составь описание профессии "%s" в одну строчку`, name)
	resp := gpt.SendRequestToGPT(question)
	return resp.Answer, resp.Error
}
