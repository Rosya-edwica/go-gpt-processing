package models

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
	Disability  []Disability
	Salary      int
}

type Disability struct {
	Id   int
	Name string
}
