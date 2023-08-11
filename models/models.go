package models

type Pair struct {
	Id int
	First string
	Second string
	IsDuplicate	bool
}

type Skill struct {
	Id int
	Name string
	IsValid bool
	Group string
}