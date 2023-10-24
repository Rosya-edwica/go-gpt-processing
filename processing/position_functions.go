package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindFunctionsForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали функции для всех профессиий"

	positions := database.GetPositionWithoutFuctions()
	for _, pos := range positions {
		functions, err := positionsGPT.GetFunctionsForPosition(pos.Name)
		checkErr(err)
		pos.Functions = functions
		database.InsertPositionFunctions(pos)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)

}
