package models

type Skill struct {
	Id            int
	Name          string
	DuplicateName string
	IsDuplicate   bool
	IsValid       bool
	GroupType     string
	SubSkills     []string
}

type Position struct {
	Id          int
	Name        string
	About       string
	Description string
	WorkPlaces  []string
	Skills      []string
	OtherNames  []string
}
