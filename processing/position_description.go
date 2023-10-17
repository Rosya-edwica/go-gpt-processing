// position_description - Собираем описание для всех профессий, у которых его нет
package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindDescriptionForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали описание для всех профессиий"
	var op = "processing.position_description"

	positions := database.GetPositionWithoutDescription()
	posCount := len(positions)
	for i, pos := range positions {
		descr, timeEx, err := positionsGPT.GetDescriptionForPosition(pos.Name)
		checkErr(err)
		pos.Description = descr
		database.UpdatePositionDescription(pos)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, posCount, pos.Name, timeEx)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)

}
