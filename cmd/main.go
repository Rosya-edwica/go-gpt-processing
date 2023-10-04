package main

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/processing"
	"os"

	"github.com/joho/godotenv"
)

var exitMessage = "Запусти программу с доп. аргументом: \n1. duplicate - чтобы запустить поиск дубликатов \n2. soft - чтобы запустить поиск soft-скиллов \n3. hard - чтобы запустить поиск hard-скиллов \n4. group - чтобы определить принадлежность к группам (навык, профессия, другое)"

func main() {
	database := initDatabase()
	detectProcessingType(database)
	database.Close()
}

func initDatabase() (database db.Database) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	database = db.Database{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		User:     os.Getenv("MYSQL_USER"),
		Name:     os.Getenv("MYSQL_NAME"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	}
	database.Connect()
	return
}

func detectProcessingType(database db.Database) {
	args := os.Args
	if len(args) == 1 {
		panic(exitMessage)
	}
	switch args[1] {
	case "duplicate":
		processing.CheckAllSkillsForDuplicates(&database)
	case "hard":
		processing.CheckAllSkillsForSoftOrHardSkill(&database, "hard")
	case "soft":
		processing.CheckAllSkillsForSoftOrHardSkill(&database, "soft")
	case "group":
		processing.CheckAllSkillsForGroupType(&database)
	case "subskills":
		processing.CollectForAllSkillsSubSkills(&database)
	case "description":
		processing.FindDescriptionForAllsPositions(&database)
	case "about":
		processing.FindAboutForAllsPositions(&database)
	case "other_names":
		processing.FindOtherNamesForAllsPositions(&database)
	case "work_places":
		processing.FindWorkPlacesForAllPositions(&database)
	case "functions":
		processing.FindFunctionsForAllPositions(&database)
	case "tests":
		processing.CollectForAllSkillsTests(&database)
	case "education":
		processing.FindEducationForAllPositions(&database)
	case "levels":
		processing.FindLevelsForAllPositions(&database)
	case "experience_salary":
		processing.FindExperienceAndSalaryForLevelPositions(&database)
	default:
		panic(exitMessage)
	}
}
