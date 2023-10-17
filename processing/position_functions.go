package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindFunctionsForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали функции для всех профессиий"
	var op = "processing.position_functions"

	positions := database.GetPositionWithoutFuctions()
	posCount := len(positions)
	for i, pos := range positions {
		functions, timeEx, err := positionsGPT.GetFunctionsForPosition(pos.Name)
		checkErr(err)
		pos.Functions = functions
		database.InsertPositionFunctions(pos)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, posCount, pos.Name, timeEx)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)

}
