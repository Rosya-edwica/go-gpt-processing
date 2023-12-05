package course

import (
	"fmt"
	"go-gpt-processing/internal/gpt"
	"regexp"
)

func GetCourseSkills(courseName string) (skills []string, err error) {
	question := fmt.Sprintf(`Какие навыки я получу после прохождения данного курса "%s"? Составь список из 20 навыков. Навыки должны быть в виде ключевых слов.`,
		courseName)
	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return nil, resp.Error
	}
	skills = parseTextToArrayByLines(resp.Answer)
	if skills == nil {
		return nil, gpt.WrongAnswerError
	}
	return

}

// Приводит текст такого формата:
// '1. Первая строка
// 2. Вторая строка
// 3. Третья строка'
// в список [Первая строка, Вторая строка, Третья строка]
func parseTextToArrayByLines(text string) (items []string) {
	reLines := regexp.MustCompile(`\d+..*`)
	reSubPointDigit := regexp.MustCompile(`\d+. `) // Вырезает такое '1. '
	lines := reLines.FindAllString(text, -1)
	if len(lines) <= 1 {
		return nil
	}

	for _, line := range lines {
		item := reSubPointDigit.ReplaceAllString(line, "")
		items = append(items, item)
	}

	return
}
