package main

import (
	"fmt"
	"gpt-skills/db"
	"gpt-skills/gpt"
	"gpt-skills/logger"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)
var duplicateFunction, softOrHardFunction, groupTypeFunction bool
var softOrHard string

var exitMessage = "Запусти программу с доп. аргументом: \n1. duplicate - чтобы запустить поиск дубликатов \n2. soft - чтобы запустить поиск soft-скиллов \n3. hard - чтобы запустить поиск hard-скиллов \n4. group - чтобы определить принадлежность к группам (навык, профессия, другое)"

func init() {
	args := os.Args
	if len(args) == 1 {
		panic(exitMessage)
	}
	if args[1] == "duplicate"{
		duplicateFunction = true
	} else if args[1] == "soft" {
		softOrHardFunction = true
		softOrHard = "soft"
	} else if args[1] == "hard" {
		softOrHardFunction = true
		softOrHard = "hard"
	} else if args[1] == "group" {
		groupTypeFunction = true
	} else {
		panic(exitMessage)
	}
}

func main() {
	err := godotenv.Load(".env")
	checkErr(err)

	database := db.Database{
		Host: os.Getenv("MYSQL_HOST"),
		Port: os.Getenv("MYSQL_PORT"),
		User: os.Getenv("MYSQL_USER"),
		Name: os.Getenv("MYSQL_NAME"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	}
	database.Connect()

	if duplicateFunction {
		CheckAllSkillsForDuplicates(&database)
	} else if softOrHardFunction {
		CheckAllSkillsForSoftOrHardSkill(&database)
	} else if groupTypeFunction {
		CheckAllSkillsForGroupType(&database)
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
		err := gpt.CheckSkillsForDuplicates(&pair)
		checkErr(err)
		if err == nil {
			database.UpdatePair(pair)
		}
	}
}

func CheckAllSkillsForSoftOrHardSkill(database *db.Database) {
	fmt.Printf("Ищем %s-skills...\n", softOrHard)
	for {
		skill := database.GetSkill(softOrHard)
		if skill.Id == 0 {
			return
		}
		err := gpt.CheckSkillIsSoftOrHard(softOrHard, &skill)
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
		err := gpt.CheckSkillsForTypeGroup(&skill)
		checkErr(err)
		if err == nil {
			database.UpdateSkillGroup(skill)
		}
		time.Sleep(2 * time.Second)
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
			panic(err)
		}
	}
}