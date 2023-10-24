package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindWorkPlacesForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали места работы для всех профессиий"

	positions := database.GetPositionWithoutWorkPlaces()
	for _, pos := range positions {
		workPlaces, err := positionsGPT.GetWorkPlacesForPosition(pos.Name)
		checkErr(err)
		pos.WorkPlaces = workPlaces
		database.UpdatePositionWorkPlaces(pos)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
