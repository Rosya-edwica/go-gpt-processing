package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
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
}

func FindDescriptionForAllsPositions(database *db.Database) {
	fmt.Println("Ищем полное описание для профессии")
	positions := database.GetPositionWithoutDescription()
	for i, pos := range positions {
		startTime := time.Now().Unix()
		descr, err := positionsGPT.GetDescriptionForPosition(pos.Name)
		if err != nil {
			fmt.Println("Ошибка для профессии:", pos.Name, err)
			Pause(10)
			continue
		}
		pos.Description = descr
		database.UpdatePositionDescription(pos)
		fmt.Printf("[%d/%d] Описание для профессии - %s (%d):\n %d seconds\n\n", i+1, len(positions), pos.Name, pos.Id, time.Now().Unix()-startTime)
		Pause(5)
	}
}

func FindOtherNamesForAllsPositions(database *db.Database) {
	fmt.Println("Ищем полное описание для профессии")
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
}
