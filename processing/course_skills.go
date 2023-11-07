// Нам нужно сопоставить навыки от GPT с теми, что у нас есть и сохранить только пересечения

package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/courseGPT"
	"go-gpt-processing/pkg/models"
	"go-gpt-processing/pkg/telegram"
)

func FindSkillsForOpeneduCourses(database *db.Database) {
	courses := database.GetOpeneduCourses()
	fmt.Println(len(courses))
	findSkillsForCourses(courses, database)
}

func FindSkillsForAllCourses(database *db.Database) {
	courses := database.GetAllCourses()
	fmt.Println(len(courses))
	findSkillsForCourses(courses, database)
}

func findSkillsForCourses(courses []models.Course, database *db.Database) {
	const SuccessMessage = "Сопоставили курсы с навыками"

	for i, course := range courses {
		skills, err := courseGPT.GetCourseSkills(course.Name)
		checkErr(err)
		existSkillsId := findGPTSkillInDatabase(database, skills)
		database.SaveCourseSkills(course.Id, existSkillsId)
		fmt.Printf("Осталось: %d/%d(%d совпадений для %s:%d)\n", i, len(courses), len(existSkillsId), course.Name, course.Id)
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
