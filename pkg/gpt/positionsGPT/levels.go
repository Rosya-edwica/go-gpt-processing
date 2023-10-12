package positionsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"go-gpt-processing/pkg/models"
	"regexp"
	"strconv"
	"strings"
)

var reLines = regexp.MustCompile(`\d+..*`)               // 1. SMM-специалист-стажер - опыт работы: 0-1 год, средняя зарплата: 30 000 рублей.
var rePositionNames = regexp.MustCompile(`\d+.* -`)      // 1. SMM-специалист-стажер -
var reExperience = regexp.MustCompile(`опыт работы:.*,`) // опыт работы: 0-1 год,
var reSalary = regexp.MustCompile(`средняя.*`)           // средняя зарплата: 30 000 рублей.
var reDigits = regexp.MustCompile(`\d+`)

var reSubPositionNames = regexp.MustCompile(`\d+. | -`)                                    // SMM-специалист-стажер
var reSubExperience = regexp.MustCompile(`опыт работы:|,`)                                 //  0-1 год
var reSubSalary = regexp.MustCompile(`средняя зарплата: |руб.*|средняя заработная плата `) // 30 000

func GetLevelsForPosition(name string) (levels []models.PositionLevel, err error) {
	question := fmt.Sprintf(`Составь уровни должности для профессии "%s". Список должен содержать только наименования уровней профессий. 
		Составь список от самого начального уровня, до самого высшего, последние самые высокие уровни должны содержать уровни топ-менеджмента и директоров,
		если данная профессия предполагает такой карьерный рост. Используй обобщенный вариант уровней профессий. 
		Не указывай в списке уровень "студент" и профессии связанные с научной деятельностью, к примеру "научный сотрудник". 
		Для каждой профессии укажи требуемый средний опыт в формате диапазона от и до. 
		Также для каждого уровня должности и требуемому опыту укажи среднюю заработную плату по России на 2021 год. 
		Зарплаты укажи в точном значении, а не диапазоны. Учитывай только российский рынок. 
		Учитывай только российские источники и ресурсы. Ответ предоставь в таком формате: "1. Уровень должности - опыт работы: срок, средняя зарплата: информация о зарплате"`,
		name,
	)
	answer, _, err := gpt.SendRequestToGPT(question)
	fmt.Println(answer)
	if err != nil {
		return []models.PositionLevel{}, err
	}
	levels = parseAnswerToPositionLevels(answer)
	if len(levels) == 0 || levels[0].Level == "" {
		return []models.PositionLevel{}, errors.New("Не удалось подобрать уровни для профессии")
	}
	return
}

func parseAnswerToPositionLevels(text string) (levels []models.PositionLevel) {
	lines := reLines.FindAllString(text, -1)
	for _, item := range lines {
		level := models.PositionLevel{
			Level:      parsePositionName(item),
			Experience: parseExperience(item),
			Salary:     parseSalary(item),
		}
		levels = append(levels, level)
	}
	return
}

func parsePositionName(text string) (name string) {
	nameText := rePositionNames.FindString(text)
	name = reSubPositionNames.ReplaceAllString(nameText, "")
	return
}

func parseExperience(text string) (experience string) {
	experienceText := reExperience.FindString(text)
	experience = reSubExperience.ReplaceAllString(experienceText, "")
	return
}

func parseSalary(text string) (num int) {
	salaryText := reSalary.FindString(text)
	salary := reSubSalary.ReplaceAllString(salaryText, "")
	salary = strings.ReplaceAll(salary, " ", "")
	salary = reDigits.FindString(salary)
	num, _ = strconv.Atoi(salary)
	return
}
