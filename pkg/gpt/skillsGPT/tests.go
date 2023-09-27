package skillsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"go-gpt-processing/pkg/models"
	"regexp"
	"strings"
)

var titleRegexp = regexp.MustCompile(`\d+. .*?\n`)
var titleRegexpSub = regexp.MustCompile(`\d+.|\n|Вопрос:`)
var choicesRegexp = regexp.MustCompile(`\w\) .*`)
var choicesRegexpSub = regexp.MustCompile(`\w\) |\n|\( ответ \)`)
var answerRegexp = regexp.MustCompile(`Ответ:.*|\w\) .*\( ответ \)`)
var answerRegexpSub = regexp.MustCompile(`Ответ: |\n|\( ответ \)|\w\) `)

func GetTestForSkill(query string) (test models.Test, err error) {
	gptQuestion := fmt.Sprintf("Для изучения навыка '%s' составь тест из 10 вопросов с вариантами ответов в такой структуре: 1. Вопрос\na) первый вариант\nb) второй вариант\nc) третий вариант\nd) четвертый вариант\nОтвет: полный вариант", query)
	answer, err := gpt.SendRequestToGPT(gptQuestion)
	if err != nil {
		fmt.Println(err)
		return models.Test{}, err
	}
	test = ParseTest(strings.TrimSpace(answer))
	return
}

func ParseTest(text string) (test models.Test) {
	listQuestionsText := strings.Split(text, "\n\n")
	for _, item := range listQuestionsText {
		question := parseQuestion(item)
		test.Questions = append(test.Questions, question)
	}
	return
}

func parseQuestion(text string) (question models.Question) {
	question.Text = titleRegexp.FindString(text)
	question.Text = titleRegexpSub.ReplaceAllString(question.Text, "")
	question.Text = strings.TrimSpace(question.Text)

	question.Choices = parseChoices(text)
	question.Answer = parseAnswer(text)
	return
}

func parseChoices(text string) (choices []string) {
	choicesFinded := choicesRegexp.FindAllString(text, -1)
	for _, item := range choicesFinded {
		item = choicesRegexpSub.ReplaceAllString(item, "")
		item = strings.TrimSpace(item)
		choices = append(choices, item)
	}
	return
}
func parseAnswer(text string) (answer string) {
	answer = answerRegexp.FindString(text)
	answer = answerRegexpSub.ReplaceAllString(answer, "")
	answer = strings.TrimSpace(answer)
	return
}
