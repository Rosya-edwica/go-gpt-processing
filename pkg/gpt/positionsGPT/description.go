package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
)

func GetDescriptionForPosition(name string) (descr string, err error) {
	question := fmt.Sprintf(`Выступи в роли карьерного эксперта и сделай описание профессии "%s". Описание должно быть написано простым языком, для людей, которые могут не знать об этой профессии. Напиши подробное описание не больше 2 абзацев. Количество символов не больше 1000`, name)
	resp := gpt.SendRequestToGPT(question)
	return resp.Answer, resp.Error
}
