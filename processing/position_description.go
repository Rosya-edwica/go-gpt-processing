// position_description - Собираем описание для всех профессий, у которых его нет
package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindDescriptionForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали описание для всех профессиий"

	positions := database.GetPositionWithoutDescription()
	for _, pos := range positions {
		descr, err := positionsGPT.GetDescriptionForPosition(pos.Name)
		checkErr(err)
		pos.Description = descr
		database.UpdatePositionDescription(pos)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)

}
