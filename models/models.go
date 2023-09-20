package models

type Pair struct {
	Id          int
	First       string
	Second      string
	IsDuplicate bool
}

type Skill struct {
	Id      int
	Name    string
	IsValid bool
	Group   string
}

type SkillForSubSkills struct {
	Id        int
	Name      string
	SubSkills []string
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
