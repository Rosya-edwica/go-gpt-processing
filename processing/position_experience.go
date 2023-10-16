package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindExperienceAndSalaryForLevelPositions(database *db.Database) {
	const SuccessMessage = "Подобрали опыт и зарплату для всех профессиий"
	var op = "processing.position_experience"

	parentIds := database.GetParentIdsForLevelsWithoutExperienceAndSalary()
	posCount := len(parentIds)
	for i, parentId := range parentIds {
		positions := database.GetPositionsLevelsWithoutExperienceAndSalaryByParentId(parentId)
		updated, timeEx, err := positionsGPT.GetLevelInfoForPosition(positions)
		if err != nil {
			fmt.Printf("%s\t ERROR:%s\n", op, err)
			Pause(30)
			continue
		}
		database.UpdatePositionsLevelExperienceAndSalary(updated, parentId)
		fmt.Printf("%s\t[%d/%d] %d (Time: %d s)\n", op, i+1, posCount, parentId, timeEx)
		Pause(3)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
