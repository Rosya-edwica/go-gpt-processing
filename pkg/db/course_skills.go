package db

import (
	"fmt"
	"go-gpt-processing/pkg/models"
)

const (
	QuerySelectOpeneduCourses = `
		SELECT course.id, course.name  FROM course
		LEFT JOIN company ON company.id = company_id
		WHERE company.id = 2077 AND course.id NOT IN (
			SELECT course_id FROM demand_to_course
		)`
	QueryInsertOpeneduSkills = `INSERT INTO demand_to_course(course_id, demand_id) VALUES(%d, %d)`
	QuerySelectSkillByName   = `SELECT id FROM demand WHERE LOWER(name) = '%s' ORDER BY is_deleted ASC LIMIT 1`
)

func (d *Database) GetOpeneduCourses() (courses []models.Course) {
	rows, err := d.Connection.Query(QuerySelectOpeneduCourses)
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

func (d *Database) SaveOpeneduSkills(courseID int, skillsId []int) {
	for _, skillID := range skillsId {
		query := fmt.Sprintf(QueryInsertOpeneduSkills, courseID, skillID)
		res, err := d.Connection.Exec(query)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		fmt.Println("Новая связь id: ", id)

	}
}

func (d *Database) GetSkillIdByName(name string) (id int) {
	query := fmt.Sprintf(QuerySelectSkillByName, name)
	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		return
	}
	return
}
