package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CheckAllSkillsForGroupType(database *db.Database) {
	const SuccessMessage = "Определили тип всех навыков"

	skills := database.GetSkillsWithoutGroup()
	for _, skill := range skills {
		groupType, err := skillsGPT.CheckSkillsForTypeGroup(skill.Name)
		checkErr(err)
		skill.GroupType = groupType
		database.UpdateSkillGroup(skill)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
