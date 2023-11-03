// Нам нужно сопоставить навыки от GPT с теми, что у нас есть и сохранить только пересечения

package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/courseGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindSkillsForOpeneduCourses(database *db.Database) {
	const SuccessMessage = "Сопоставили курсы открытого образования с навыками"

	courses := database.GetOpeneduCourses()
	for i, course := range courses {
		skills, err := courseGPT.GetCourseSkills(course.Name)
		checkErr(err)
		existSkillsId := findGPTSkillInDatabase(database, skills)
		database.SaveOpeneduSkills(course.Id, existSkillsId)
		fmt.Printf("Осталось: %d/%d(%d совпадений для %s)\n", i, len(courses), len(existSkillsId), course.Name)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}

func findGPTSkillInDatabase(database *db.Database, skills []string) (ids []int) {
	for _, skill := range skills {
		id := database.GetSkillIdByName(skill)
		if id != 0 {
			ids = append(ids, id)
		}
	}
	return
}
