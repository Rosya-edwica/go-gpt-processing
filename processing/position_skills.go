package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/telegram"
)

func FindSkillsForPositions(database *db.Database) {
	profAreaas := database.GetProfAreaList()
	for _, area := range profAreaas {
		FindSkillsInProfArea(database, area)
	}
}

func FindSkillsInProfArea(database *db.Database, area string) {
	positions := database.GetPositionsByProfArea(area)
	posCount := len(positions)
	if posCount == 0 {
		return
	}
	for _, pos := range positions {
		skills, err := positionsGPT.GetSkillsForPosition(pos.Name, pos.ProfArea)
		checkErr(err)
		pos.Skills = skills
		database.SavePositionSkills(pos)

	}
	SuccessMessage := fmt.Sprintf("Подобрали навыки для ненулевых профессий из профобласти - %s\nКоличество профессий:%d", area, posCount)
	telegram.SuccessMessageMailing(SuccessMessage)
}
