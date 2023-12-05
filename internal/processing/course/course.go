// Нам нужно сопоставить навыки от GPT с теми, что у нас есть и сохранить только пересечения

package course

import (
	"fmt"
	"go-gpt-processing/internal/models"
	"go-gpt-processing/pkg/gpt/courseGPT"
	"go-gpt-processing/pkg/telegram"
	"go-gpt-processing/tools"

	rep "go-gpt-processing/internal/repositories/course"

	"github.com/go-faster/errors"
)

func FindSkillsForOpeneduCourses(r *rep.Repository) error {
	courses, err := r.GetOpeneduCoursesWithoutGPTSkills()
	if err != nil {
		return errors.Wrap(err, "processing openedu skills")
	}
	findSkillsForCourses(courses, r)
	return err
}

func FindSkillsForAllCourses(r *rep.Repository) error {
	courses, err := r.GetAllCoursesWithoutGPTSkills()
	if err != nil {

		return errors.Wrap(err, "processing openedu skills")
	}
	fmt.Println(len(courses))
	findSkillsForCourses(courses, r)
	return nil
}

func findSkillsForCourses(courses []models.Course, r *rep.Repository) {
	const SuccessMessage = "Сопоставили курсы с навыками"

	for i, course := range courses {
		skills, err := courseGPT.GetCourseSkills(course.Name)
		tools.CheckErr(err)
		existSkillsId := findGPTSkillInDatabase(r, skills)
		r.ConnectCourseWithSkills(course.Id, existSkillsId)
		fmt.Printf("Осталось: %d/%d(%d совпадений для %s:%d)\n", i, len(courses), len(existSkillsId), course.Name, course.Id)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}

func findGPTSkillInDatabase(r *rep.Repository, skills []string) (ids []int64) {
	for _, skill := range skills {
		id, err := r.GetSkillIdByName(skill)
		tools.CheckErr(err)
		if id != 0 {
			ids = append(ids, id)
		}
	}
	return
}
