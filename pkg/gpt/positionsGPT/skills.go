package positionsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

func GetSkillsForPosition(name string) (skills []string, err error) {
	question := fmt.Sprintf(`
	Составь список из 90 профессиональных навыков и знаний для профессии %s.  Пиши в строчку. Не используй нумерацию. 
	Не используй в описании такие слова как: знание, умение, владение, работа с, понимание. В качестве разделителя используй знак ,`,
		name)
	answer, err := gpt.SendRequestToGPT(question)
	skills = strings.Split(answer, ",")

	if len(skills) <= 1 {
		return nil, errors.New(fmt.Sprintf("Не удалось поделить ответ по запятым: %s", answer))
	}
	if answer == "" {
		return nil, errors.New(fmt.Sprintf("Пустое ответ для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", answer, name))
	}
	return
}
