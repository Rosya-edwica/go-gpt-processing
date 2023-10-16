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
	var op = "processing.position_skills"

	positions := database.GetPositionsByProfArea(area)
	posCount := len(positions)
	if posCount == 0 {
		return
	}
	for i, pos := range positions {
		skills, timeEx, err := positionsGPT.GetSkillsForPosition(pos.Name, pos.ProfArea)
		if err != nil {
			fmt.Printf("%s\t ERROR:%s\n", op, err)
			Pause(30)
			continue
		}
		if len(skills) < 10 {
			continue
		}
		pos.Skills = skills
		database.SavePositionSkills(pos)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, posCount, pos.Name, timeEx)

	}
	SuccessMessage := fmt.Sprintf("Подобрали навыки для ненулевых профессий из профобласти - %s\nКоличество профессий:%d", area, posCount)
	telegram.SuccessMessageMailing(SuccessMessage)
}
