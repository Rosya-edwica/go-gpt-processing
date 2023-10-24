package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindOtherNamesForAllsPositions(database *db.Database) {
	const SuccessMessage = "Подобрали другие наименования для всех профессиий"

	positions := database.GetPositionWithoutOtherNames()
	for _, pos := range positions {
		otherNames, err := positionsGPT.GetOtherNamesForPosition(pos.Name)
		checkErr(err)
		pos.OtherNames = otherNames
		database.UpdatePositionOtherNames(pos)
		Pause(3)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
