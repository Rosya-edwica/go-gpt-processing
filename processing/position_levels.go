package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindLevelsForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали уровни для всех профессиий"
	var op = "processing.position_levels"

	positions := database.GetPositionsWithoutLevels()
	posCount := len(positions)
	for i, pos := range positions {
		levels, timeEx, err := positionsGPT.GetLevelsForPosition(pos.Name)
		checkErr(err)
		pos.Levels = levels
		database.InsertPositionLevels(pos)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, posCount, pos.Name, timeEx)
		Pause(3)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
