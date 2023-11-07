package db

import (
	"fmt"
	"go-gpt-processing/pkg/models"
	"strings"
)

const (
	QuerySelectAllCourses = `
		SELECT id, name 
		FROM course 
		WHERE is_public = 1
		AND id NOT IN (
			SELECT course_id FROM demand_to_course WHERE is_chatgpt IS TRUE
		)
		`

	QuerySelectOpeneduCourses = `
		SELECT course.id, course.name  FROM course
		LEFT JOIN company ON company.id = company_id
		WHERE company.id = 2077 AND course.id NOT IN (
			SELECT course_id FROM demand_to_course
		)`
	QueryInsertSkills      = `INSERT INTO demand_to_course(course_id, demand_id, is_chatgpt) VALUES(%d, %d, true)`
	QuerySelectSkillByName = `SELECT id FROM demand WHERE LOWER(name) = '%s' ORDER BY is_deleted ASC LIMIT 1`
)

func (d *Database) GetOpeneduCourses() (courses []models.Course) {
	return d.GetCoursesByQuery(QuerySelectOpeneduCourses)
}

func (d *Database) GetAllCourses() (courses []models.Course) {
	return d.GetCoursesByQuery(QuerySelectAllCourses)
}

func (d *Database) GetCoursesByQuery(query string) (courses []models.Course) {
	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		courses = append(courses, models.Course{
			Id:   id,
			Name: name,
		})
	}
	return
}

func (d *Database) SaveCourseSkills(courseID int, skillsId []int) {
	for _, skillID := range skillsId {
		query := fmt.Sprintf(QueryInsertSkills, courseID, skillID)
		res, err := d.Connection.Exec(query)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		fmt.Println("Новая связь id: ", id)

	}
}

func (d *Database) GetSkillIdByName(name string) (id int) {
	query := fmt.Sprintf(QuerySelectSkillByName, strings.ReplaceAll(name, "'", ""))
	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		return
	}
	return
}
