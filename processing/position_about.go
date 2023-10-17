// position_about - Собираем краткое описание для всех профессий, у которых его нет
package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindAboutForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали описание для всех профессиий"
	var op = "processing.position_about"

	positions := database.GetPositionWithoutAbout()
	posCount := len(positions)
	for i, pos := range positions {
		about, timeEx, err := positionsGPT.GetAboutForPosition(pos.Name)
		checkErr(err)
		pos.About = about
		database.UpdatePositionAbout(pos)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, posCount, pos.Name, timeEx)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
