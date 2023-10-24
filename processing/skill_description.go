package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CollectDescriptionForAllSkills(database *db.Database) {
	const SuccessMessage = "Подобрали описание для всех навыков"

	skills := database.GetNullableSkillsInColumn("description")
	for _, skill := range skills {
		description, err := skillsGPT.GetDescriptionForSkill(skill.Name)
		checkErr(err)
		skill.Description = description
		database.UpdateSkillColumn(skill.Id, "description", skill.Description)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
