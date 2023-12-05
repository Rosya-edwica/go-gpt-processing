package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindDisabilityForAllPositions(database *db.Database) {
	const SuccessMessage = "Определили инвалидность для всех профессий"

	positions := database.GetPositionsWithoutDisability()
	disabilities := database.GetDisablities()
	for _, pos := range positions {
		posDisability, err := positionsGPT.GetDisabilityForPosition(pos.Name, disabilities)
		checkErr(err)
		pos.Disability = posDisability
		fmt.Printf("Инвалидность для: '%s' -> %#v\n", pos.Name, posDisability)
		database.UpdatePositionDisability(pos)
		Pause(3)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
