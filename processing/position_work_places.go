package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindWorkPlacesForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали места работы для всех профессиий"
	var op = "processing.position_work_places"

	positions := database.GetPositionWithoutWorkPlaces()
	posCount := len(positions)
	for i, pos := range positions {
		workPlaces, timeEx, err := positionsGPT.GetWorkPlacesForPosition(pos.Name)
		if err != nil {
			fmt.Printf("%s\t ERROR:%s\n", op, err)
			Pause(30)
			continue
		}
		pos.WorkPlaces = workPlaces
		database.UpdatePositionWorkPlaces(pos)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, posCount, pos.Name, timeEx)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
