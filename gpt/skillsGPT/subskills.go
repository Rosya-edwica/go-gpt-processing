package skillsGPT

import (
	"errors"
	"fmt"
	"gpt-skills/gpt"
	"regexp"
)

func GetSubSkills(query string) (skills []string, err error) {
	question := fmt.Sprintf(`
	какими hard-skills нужно обладать, чтобы изучить "%s" - сократи ответ до перечня навыков
	`, query)

	answer, err := gpt.SendRequestToGPT(question)
	if err != nil {
		return []string{}, errors.New(fmt.Sprintf("ОШИБКА: %s", err.Error()))
	}
	reLines := regexp.MustCompile(`\d+. .*?\n`)
	reDigits := regexp.MustCompile(`\d+. |\n`)
	lines := reLines.FindAllString(answer+"\n", -1)

	for _, line := range lines {
		skill := reDigits.ReplaceAllString(line, "")
		skills = append(skills, skill)
	}
	return
}
