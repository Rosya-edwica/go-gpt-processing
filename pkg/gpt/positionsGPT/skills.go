package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"regexp"
	"strings"
)

var rePointDigit = regexp.MustCompile(`\d+. `)

func GetSkillsForPosition(name string, profarea string) (skills []string, err error) {
	question := fmt.Sprintf(`
	Составь список из 20 навыков и знаний для профессии  "%s". 
	Пиши коротко, не более четырех слов. Не указывай банальные знания и навыки, пиши только Hard Skills и не указывай Soft Skills. 
	Не используй в описании такие слова как: знание, умение, владение, работа с, понимание, использование. 
	Навыки должны относиться к профобласти "%s"`,
		name, profarea)
	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return nil, resp.Error
	}

	skills = reLines.FindAllString(resp.Answer, -1)
	if len(skills) <= 1 {
		return nil, gpt.WrongAnswerError
	}

	var skillsWithoutDigitPoint []string
	for _, i := range skills {
		skill := rePointDigit.ReplaceAllString(i, "")
		skill = strings.ReplaceAll(skill, ".", "")
		skillsWithoutDigitPoint = append(skillsWithoutDigitPoint, strings.TrimSpace(skill))
	}
	return skillsWithoutDigitPoint, nil
}
