package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
	"strings"
	"time"
)

func FindAboutForAllsPositions(database *db.Database) {
	fmt.Println("Ищем описание для профессии")
	positions := database.GetPositionWithoutAbout()
	for i, pos := range positions {
		about, err := positionsGPT.GetAboutForPosition(pos.Name)
		if err != nil {
			Pause(120)
			about, err = positionsGPT.GetAboutForPosition(pos.Name)
			if err != nil {
				continue
			}
		}
		pos.About = about
		database.UpdatePositionAbout(pos)
		fmt.Printf("[%d/%d] Описание для профессии - %s (%d):\n %s\n\n", i+1, len(positions), pos.Name, pos.Id, pos.About)
		Pause(5)
	}
	telegram.SuccessMessageMailing("Поиск краткого описания для профессий завершился успешно")

}

func FindDescriptionForAllsPositions(database *db.Database) {
	fmt.Println("Ищем полное описание для профессии")
	positions := database.GetPositionWithoutDescription()
	for i, pos := range positions {
		startTime := time.Now().Unix()
		descr, err := positionsGPT.GetDescriptionForPosition(pos.Name)
		if err != nil {
			if strings.Contains(err.Error(), "status code: 429") {
				checkErr(err)
			}
			fmt.Println("Ошибка для профессии:", pos.Name, err)
			Pause(10)
			continue
		}
		pos.Description = descr
		database.UpdatePositionDescription(pos)
		fmt.Printf("[%d/%d] Описание для профессии - %s (%d):\n %d seconds\n\n", i+1, len(positions), pos.Name, pos.Id, time.Now().Unix()-startTime)
		Pause(5)
	}
	telegram.SuccessMessageMailing("Поиск подробного описания для профессий завершился успешно")

}

func FindOtherNamesForAllsPositions(database *db.Database) {
	fmt.Println("Ищем другие наименования для профессии")
	positions := database.GetPositionWithoutOtherNames()
	for i, pos := range positions {
		otherNames, err := positionsGPT.GetOtherNamesForPosition(pos.Name)
		if err != nil {
			Pause(120)
			otherNames, err = positionsGPT.GetOtherNamesForPosition(pos.Name)
			if err != nil {
				continue
			}
		}
		pos.OtherNames = otherNames
		database.UpdatePositionOtherNames(pos)
		fmt.Printf("[%d/%d] Другие написания для профессии - %s (%d):\n %s\n\n", i+1, len(positions), pos.Name, pos.Id, pos.OtherNames)
		Pause(5)
	}
	telegram.SuccessMessageMailing("Поиск других наименований для профессий завершился успешно")

}

func FindWorkPlacesForAllPositions(database *db.Database) {
	fmt.Println("Ищем места работ для профессии")
	positions := database.GetPositionWithoutWorkPlaces()
	for i, pos := range positions {
		workPlaces, err := positionsGPT.GetWorkPlacesForPosition(pos.Name)
		if err != nil {
			Pause(120)
			workPlaces, err = positionsGPT.GetWorkPlacesForPosition(pos.Name)
			if err != nil {
				continue
			}
		}
		pos.WorkPlaces = workPlaces
		database.UpdatePositionWorkPlaces(pos)
		fmt.Printf("[%d/%d] Места работы для профессии - %s (%d):\n %s\n\n", i+1, len(positions), pos.Name, pos.Id, pos.WorkPlaces)
		Pause(5)
	}
	telegram.SuccessMessageMailing("Поиск рабочих мест для профессий завершился успешно")

}

func FindFunctionsForAllPositions(database *db.Database) {
	fmt.Println("Ищем функции для профессии")
	positions := database.GetPositionWithoutFuctions()
	for i, pos := range positions {
		startTime := time.Now().Unix()
		functions, err := positionsGPT.GetFunctionsForPosition(pos.Name)
		checkErr(err)
		pos.Functions = functions
		database.InsertPositionFunctions(pos)
		fmt.Printf("[%d/%d] Функции для профессии - %s (%d):\n %d seconds.\n\n", i+1, len(positions), pos.Name, pos.Id, time.Now().Unix()-startTime)
		Pause(5)
	}
	telegram.SuccessMessageMailing("Поиск функций для профессий завершился успешно")
}

func FindEducationForAllPositions(database *db.Database) {
	fmt.Println("Ищем уровни образования для профессий")

	for {
		posCount := database.CountPositionsWithoutEducation()
		pos := database.GetOnePositionWithoutEducation()
		if pos.Id == 0 {
			break
		}
		startTime := time.Now().Unix()
		education, err := positionsGPT.GetEducationForPosition(pos.Name)
		if err != nil {
			if err.Error() == "Не получилось распарсить результат" {
				fmt.Println("Пустой результат для профессии:", pos.Name, pos.Id)
				pos.Education = []string{""}
			} else {
				checkErr(err)
			}
		}
		pos.Education = education
		database.UpdatePositionEducation(pos)
		fmt.Printf("[Осталось: %d] Образование для профессии - %s (id:%d):%s\n %d seconds.\n\n", posCount, pos.Name, pos.Id, pos.Education, time.Now().Unix()-startTime)
		Pause(3)
	}
	telegram.SuccessMessageMailing("Поиск уровней образования для профессий завершился успешно")

}

func FindLevelsForAllPositions(database *db.Database) {
	fmt.Println("Подбираем уровни для профессий")
	positions := database.GetPositionsWithoutLevels()

	posCount := len(positions)
	for i, pos := range positions {
		startTime := time.Now().Unix()
		levels, err := positionsGPT.GetLevelsForPosition(pos)
		checkErr(err)
		pos.Levels = levels
		fmt.Println(pos)
		database.InsertPositionLevels(pos)
		fmt.Printf("[Осталось: %d] Уровни для профессии - %s (id:%d):%s\n %d seconds.\n\n", posCount-(i+1), pos.Name, pos.Id, pos.Levels[0].Level, time.Now().Unix()-startTime)
		Pause(3)

	}
}
