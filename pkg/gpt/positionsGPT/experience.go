// Здесь мы будем обрабатывать профессии не-нулевого уровня, у которых нет опыта или зарплаты

package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"go-gpt-processing/pkg/models"
	"regexp"
	"strconv"
	"strings"
)

var reLevelLine = regexp.MustCompile(`.* - .* -.*`)           // Ассистент отдела - 30000 рублей - 1-2 года;
var reLevelExperience = regexp.MustCompile(`руб.*`)           // рублей - 2-3 года;
var reSubLevelExperience = regexp.MustCompile(`руб.* -|;|\.`) // от 1 до 3 лет
var reLevelSalary = regexp.MustCompile(`.* руб`)              // - 60000 руб
var reSubLevelSalary = regexp.MustCompile(`руб|-`)            //  35000 рублей

type LevelInfo struct {
	Experience string
	Salary     int
}

func GetLevelInfoForPosition(positions []models.Position) (updated []models.Position, err error) {
	names := allPositionsNameToString(positions)
	question := fmt.Sprintf(`Напиши зарплату в рублях и опыт для профессий ниже: %s
	Зарплата должна расти в зависимости от уровня профессии, используй данные на 2021 год, используй только российские источники. 
	Ответ запиши в формате: [профессия - зарплата - опыт]. Пояснений давать не нужно. 
	Зарплату необходимо указать в абсолютном значении без диапазона. Опыт - в диапазоне.`, names)
	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return []models.Position{}, resp.Error
	}

	lines := reLines.FindAllString(resp.Answer, -1)
	if len(lines) != len(positions) {
		return []models.Position{}, gpt.WrongAnswerError
	}
	for i, pos := range positions {
		levelInfo := parseAnswerToLevelInfo(lines[i])
		pos.Experience = levelInfo.Experience
		pos.Salary = levelInfo.Salary
		updated = append(updated, pos)
	}

	return
}

func allPositionsNameToString(positions []models.Position) (names string) {
	var namesList []string
	for _, item := range positions {
		namesList = append(namesList, item.Name)
	}
	names = "\n" + strings.Join(namesList, "\n")
	return
}

func parseAnswerToLevelInfo(text string) (info LevelInfo) {
	experience := reLevelExperience.FindString(text)
	experience = reSubLevelExperience.ReplaceAllString(experience, "")
	experience = strings.TrimSpace(experience)

	salary := reLevelSalary.FindString(text)
	salary = reSubLevelSalary.ReplaceAllString(salary, "")
	salary = strings.ReplaceAll(salary, " ", "")
	salary = reDigits.FindString(salary)
	num, _ := strconv.Atoi(salary)

	return LevelInfo{
		Experience: experience,
		Salary:     num,
	}
}
