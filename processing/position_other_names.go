package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindOtherNamesForAllsPositions(database *db.Database) {
	const SuccessMessage = "Подобрали другие наименования для всех профессиий"
	var op = "processing.position_other_names"

	positions := database.GetPositionWithoutOtherNames()
	posCount := len(positions)
	for i, pos := range positions {
		otherNames, timeEx, err := positionsGPT.GetOtherNamesForPosition(pos.Name)
		checkErr(err)
		pos.OtherNames = otherNames
		database.UpdatePositionOtherNames(pos)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, posCount, pos.Name, timeEx)
		Pause(3)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
