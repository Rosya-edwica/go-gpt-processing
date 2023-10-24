package skillsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"regexp"
)

func GetSubSkills(query string) (skills []string, err error) {
	question := fmt.Sprintf(`
	какими hard-skills нужно обладать, чтобы изучить "%s" - сократи ответ до перечня навыков
	`, query)

	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return nil, resp.Error
	}
	reLines := regexp.MustCompile(`\d+. .*?\n`)
	reDigits := regexp.MustCompile(`\d+. |\n`)
	lines := reLines.FindAllString(resp.Answer+"\n", -1)

	for _, line := range lines {
		skill := reDigits.ReplaceAllString(line, "")
		skills = append(skills, skill)
	}
	return
}
