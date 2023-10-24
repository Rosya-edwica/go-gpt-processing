package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindExperienceAndSalaryForLevelPositions(database *db.Database) {
	const SuccessMessage = "Подобрали опыт и зарплату для всех профессиий"

	parentIds := database.GetParentIdsForLevelsWithoutExperienceAndSalary()
	for _, parentId := range parentIds {
		positions := database.GetPositionsLevelsWithoutExperienceAndSalaryByParentId(parentId)
		updated, err := positionsGPT.GetLevelInfoForPosition(positions)
		checkErr(err)
		database.UpdatePositionsLevelExperienceAndSalary(updated, parentId)
		Pause(3)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
