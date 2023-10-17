package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CollectDescriptionForAllSkills(database *db.Database) {
	const SuccessMessage = "Подобрали описание для всех навыков"
	var op = "processing.skill_description"

	skills := database.GetNullableSkillsInColumn("description")
	skillsCount := len(skills)
	for i, skill := range skills {
		description, exTime, err := skillsGPT.GetDescriptionForSkill(skill.Name)
		checkErr(err)
		skill.Description = description
		database.UpdateSkillColumn(skill.Id, "description", skill.Description)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, skillsCount, skill.Name, exTime)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
