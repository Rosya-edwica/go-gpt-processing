package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
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
		descr, err := positionsGPT.GetDescriptionForPosition(pos.Name)
		if err != nil {
			Pause(120)
			descr, err = positionsGPT.GetDescriptionForPosition(pos.Name)
			if err != nil {
				continue
			}
		}
		pos.Description = descr
		database.UpdatePositionDescription(pos)
		fmt.Printf("[%d/%d] Описание для профессии - %s (%d):\n %s\n\n", i+1, len(positions), pos.Name, pos.Id, pos.Description)
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
