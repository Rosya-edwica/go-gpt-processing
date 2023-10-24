package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindLevelsForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали уровни для всех профессиий"

	positions := database.GetPositionsWithoutLevels()
	for _, pos := range positions {
		levels, err := positionsGPT.GetLevelsForPosition(pos.Name)
		checkErr(err)
		pos.Levels = levels
		database.InsertPositionLevels(pos)
		Pause(3)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
