package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindExtraSkillsForAllPositions(database *db.Database) {
	const SuccessMessage = "Подобрали навыки для профессий"

	positions := database.GetPositionsWithoutGPTSkills()
	fmt.Println(len(positions), "count_pos")
	for i, pos := range positions {
		skills, err := positionsGPT.GetExtraSkillsForPosition(pos.Name, pos.Skills)
		checkErr(err)
		database.SaveExtraSkills(pos, skills)
		Pause(3)
		fmt.Printf("Осталось: %d/%d\n ", i+1, len(positions))
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
