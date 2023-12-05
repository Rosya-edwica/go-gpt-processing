package models

import "go-gpt-processing/internal/entities"

type Skill struct {
	Id            int
	Name          string
	DuplicateName string
	IsDuplicate   bool
	IsValid       bool
	GroupType     string
	SubSkills     []string
	Description   string
}

type Question struct {
	Text    string
	Choices []string
	Answer  string
}

type Test struct {
	Questions []Question
}

func NewSkills(rawSkills []entities.Skill) (skills []Skill) {
	for _, skill := range rawSkills {
		skills = append(skills, Skill{
			Id:            skill.Id,
			Name:          skill.Name,
			GroupType:     skill.GroupType,
			DuplicateName: skill.DuplicateName,
			IsDuplicate:   skill.IsDuplicate,
			IsValid:       skill.IsValid,
		})
	}
	return
}
