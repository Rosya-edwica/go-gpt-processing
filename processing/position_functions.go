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
		if err != nil {
			fmt.Printf("%s\t ERROR:%s\n", op, err)
			Pause(30)
			continue
		}
		pos.Functions = functions
		database.InsertPositionFunctions(pos)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, posCount, pos.Name, timeEx)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)

}
