// position_about - Собираем краткое описание для всех профессий, у которых его нет
package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindAboutForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали описание для всех профессиий"

	positions := database.GetPositionWithoutAbout()
	for _, pos := range positions {
		about, err := positionsGPT.GetAboutForPosition(pos.Name)
		checkErr(err)
		pos.About = about
		database.UpdatePositionAbout(pos)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
