package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CheckAllSkillsForGroupType(database *db.Database) {
	const SuccessMessage = "Определили тип всех навыков"
	var op = "processing.skill_group_type"

	skills := database.GetSkillsWithoutGroup()
	skillsCount := len(skills)
	for i, skill := range skills {
		groupType, exTime, err := skillsGPT.CheckSkillsForTypeGroup(skill.Name)
		checkErr(err)
		skill.GroupType = groupType
		database.UpdateSkillGroup(skill)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, skillsCount, skill.Name, exTime)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
