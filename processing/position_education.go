package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindEducationForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали другие наименования для всех профессиий"
	var op = "processing.position_education"

	for {
		posCount := database.CountPositionsWithoutEducation()
		pos := database.GetOnePositionWithoutEducation()
		if pos.Id == 0 {
			break
		}
		education, timeEx, err := positionsGPT.GetEducationForPosition(pos.Name)
		if err != nil {

		}
		pos.Education = education
		database.UpdatePositionEducation(pos)
		fmt.Printf("%s\t[Осталось: %d] %s (Time: %d s)\n", op, posCount, pos.Name, timeEx)
		Pause(3)
	}
	telegram.SuccessMessageMailing("Поиск уровней образования для профессий завершился успешно")

}
