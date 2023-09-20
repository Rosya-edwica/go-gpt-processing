package main

import (
	"fmt"
	"gpt-skills/db"
	"gpt-skills/gpt/positionsGPT"
	"gpt-skills/gpt/skillsGPT"
	"gpt-skills/logger"
	"gpt-skills/telegram"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var duplicateFunction, softOrHardFunction, groupTypeFunction bool
var softOrHard string

var exitMessage = "Запусти программу с доп. аргументом: \n1. duplicate - чтобы запустить поиск дубликатов \n2. soft - чтобы запустить поиск soft-скиллов \n3. hard - чтобы запустить поиск hard-скиллов \n4. group - чтобы определить принадлежность к группам (навык, профессия, другое)"

func main() {
	args := os.Args
	if len(args) == 1 {
		panic(exitMessage)
	}
	err := godotenv.Load(".env")
	checkErr(err)

	database := db.Database{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		User:     os.Getenv("MYSQL_USER"),
		Name:     os.Getenv("MYSQL_NAME"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	}
	database.Connect()
	switch args[1] {
	case "duplicate":
		CheckAllSkillsForDuplicates(&database)
	case "hard":
		CheckAllSkillsForSoftOrHardSkill(&database, "hard")
	case "soft":
		CheckAllSkillsForSoftOrHardSkill(&database, "soft")
	case "group":
		CheckAllSkillsForGroupType(&database)
	case "subskills":
		CollectForAllSkillsSubSkills(&database)
	case "description":
		FindDescriptionForAllPositions(&database)
	case "about":
		FindAboutForAllPositions(&database)
	case "other_names":
		FindOtherNamesForAllPositions(&database)
	case "work_places":
		FindWorkPlacesForAllPositions(&database)
	case "skills":
		FinSkillsForAllPositions(&database)
	default:
		panic(exitMessage)
	}
	database.Close()
}

func CheckAllSkillsForDuplicates(database *db.Database) {
	fmt.Println("Ищем дубликаты...")
	for {
		pair := database.GetSkillsPair()
		if pair.Id == 0 {
			return
		}
		err := skillsGPT.CheckSkillsForDuplicates(&pair)
		checkErr(err)
		if err == nil {
			database.UpdatePair(pair)
		}
	}
}

func CheckAllSkillsForSoftOrHardSkill(database *db.Database, softOrHard string) {
	fmt.Printf("Ищем %s-skills...\n", softOrHard)
	for {
		skill := database.GetSkill(softOrHard)
		if skill.Id == 0 {
			return
		}
		err := skillsGPT.CheckSkillIsSoftOrHard(softOrHard, &skill)
		checkErr(err)
		if err == nil {
			database.UpdateSkill(softOrHard, skill)
		}
	}
}

func CheckAllSkillsForGroupType(database *db.Database) {
	fmt.Println("Пытаемся определить группу для навыка")
	for {
		skill := database.GetSkillWithoutGroup()
		if skill.Id == 0 {
			return
		}
		err := skillsGPT.CheckSkillsForTypeGroup(&skill)
		checkErr(err)
		if err == nil {
			database.UpdateSkillGroup(skill)
		}
		time.Sleep(5 * time.Second)
	}
}

func CollectForAllSkillsSubSkills(database *db.Database) {
	fmt.Println("Ищем поднавыки")
	skills := database.GetSkills()
	for i, skill := range skills {
		subskills, err := skillsGPT.GetSubSkills(skill.Name)
		checkErr(err)
		if len(subskills) != 0 {
			skill.SubSkills = subskills
			database.SaveSubskills(skill)
		}
		fmt.Printf("[%d/%d] Поднавыки для навыка - %s:\n %s\n\n", i+1, len(skills), skill.Name, strings.Join(skill.SubSkills, "|"))

	}
}

func FindAboutForAllPositions(database *db.Database) {
	fmt.Println("Ищем описание для профессии")
	positions := database.GetPositionWithoutAbout()
	for i, pos := range positions {
		about, err := positionsGPT.GetAboutForPosition(pos.Name)
		checkErr(err)
		pos.About = about
		database.UpdatePositionAbout(pos)
		fmt.Printf("[%d/%d] Описание для профессии - %s (%d):\n %s\n\n", i+1, len(positions), pos.Name, pos.Id, pos.About)
		time.Sleep(time.Second * 10)

	}
}

func FindDescriptionForAllPositions(database *db.Database) {
	fmt.Println("Ищем полное описание для профессии")
	positions := database.GetPositionWithoutDescription()
	for i, pos := range positions {
		descr, err := positionsGPT.GetDescriptionForPosition(pos.Name)
		checkErr(err)
		pos.Description = descr
		database.UpdatePositionDescription(pos)
		fmt.Printf("[%d/%d] Описание для профессии - %s (%d):\n %s\n\n", i+1, len(positions), pos.Name, pos.Id, pos.Description)
		time.Sleep(time.Second * 10)

	}
}

func FindOtherNamesForAllPositions(database *db.Database) {
	fmt.Println("Ищем полное описание для профессии")
	positions := database.GetPositionWithoutOtherNames()
	for i, pos := range positions {
		otherNames, err := positionsGPT.GetOtherNamesForPosition(pos.Name)
		checkErr(err)
		pos.OtherNames = otherNames
		database.UpdatePositionOtherNames(pos)
		fmt.Printf("[%d/%d] Другие написания для профессии - %s (%d):\n %s\n\n", i+1, len(positions), pos.Name, pos.Id, pos.OtherNames)
		time.Sleep(time.Second * 10)

	}
}

func FindWorkPlacesForAllPositions(database *db.Database) {
	fmt.Println("Ищем места работ для профессии")
	positions := database.GetPositionWithoutWorkPlaces()
	for i, pos := range positions {
		workPlaces, err := positionsGPT.GetWorkPlacesForPosition(pos.Name)
		checkErr(err)
		pos.WorkPlaces = workPlaces
		database.UpdatePositionWorkPlaces(pos)
		fmt.Printf("[%d/%d] Места работы для профессии - %s (%d):\n %s\n\n", i+1, len(positions), pos.Name, pos.Id, pos.WorkPlaces)
		time.Sleep(time.Second * 10)

	}
}

func FinSkillsForAllPositions(database *db.Database) {
	fmt.Println("Ищем навыки для профессии")
	positions := database.GetPositionWithoutSkills()
	for i, pos := range positions {
		skills, err := positionsGPT.GetSkillsForPosition(pos.Name)
		checkErr(err)
		pos.Skills = skills
		// database.UpdatePositionSkills(pos)
		fmt.Printf("[%d/%d] Навыки для профессии - %s (%d):\n %s\n\n", i+1, len(positions), pos.Name, pos.Id, pos.Skills)
		time.Sleep(time.Second * 10)

	}
}

func checkErr(err error) {
	if err != nil {
		if strings.HasPrefix(err.Error(), "Неправильный ответ") {
			logger.Log.Printf("ОШИБКА: %s", err)
			return
		} else if strings.Contains(err.Error(), "context deadline exceeded") {
			logger.Log.Printf("ОШИБКА: %s", err)
			time.Sleep(time.Second * 10)
		} else if strings.Contains(err.Error(), "status code: 503") {
			logger.Log.Printf("ОШИБКА: %s", err)
			time.Sleep(time.Second * 10)
		} else {
			telegram.Mailing(err.Error())
			panic(err)
		}
	}
}
