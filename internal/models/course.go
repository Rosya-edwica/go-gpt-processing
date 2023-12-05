package models

import "go-gpt-processing/internal/entities"

type Course struct {
	Id   int
	Name string
}

func NewCourses(rawCourses []entities.Course) (courses []Course) {
	for _, course := range rawCourses {
		courses = append(courses, Course{
			Id:   course.Id,
			Name: course.Name,
		})
	}
	return
}
