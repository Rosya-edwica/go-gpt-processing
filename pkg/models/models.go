package models

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
type PositionLevel struct {
	Level      string
	Experience string
	Salary     int
}

type Position struct {
	Id          int
	Name        string
	ProfArea    string
	About       string
	Description string
	WorkPlaces  []string
	Skills      []string
	OtherNames  []string
	Functions   []string
	Education   []string
	Levels      []PositionLevel
	Experience  string
	Salary      int
}

type Question struct {
	Text    string
	Choices []string
	Answer  string
}

type Test struct {
	Questions []Question
}
