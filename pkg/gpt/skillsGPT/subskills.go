package skillsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"regexp"
)

func GetSubSkills(query string) (skills []string, exTime int64, err error) {
	question := fmt.Sprintf(`
	какими hard-skills нужно обладать, чтобы изучить "%s" - сократи ответ до перечня навыков
	`, query)

	answer, exTime, err := gpt.SendRequestToGPT(question)
	if err != nil {
		return []string{}, exTime, errors.New(fmt.Sprintf("ОШИБКА: %s", err.Error()))
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
