package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindEducationForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали другие наименования для всех профессиий"

	for {
		pos := database.GetOnePositionWithoutEducation()
		if pos.Id == 0 {
			break
		}
		education, err := positionsGPT.GetEducationForPosition(pos.Name)
		checkErr(err)
		pos.Education = education
		database.UpdatePositionEducation(pos)
		Pause(3)
	}
	telegram.SuccessMessageMailing("Поиск уровней образования для профессий завершился успешно")

}
