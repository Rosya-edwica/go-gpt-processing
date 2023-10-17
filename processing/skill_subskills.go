package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CollectForAllSkillsSubSkills(database *db.Database) {
	const SuccessMessage = "Подобрали поднавыки для всех навыков"
	var op = "processing.skill_subskills"

	skills := database.GetSkills()
	skillsCount := len(skills)

	for i, skill := range skills {
		subskills, exTime, err := skillsGPT.GetSubSkills(skill.Name)
		checkErr(err)
		if len(subskills) != 0 {
			skill.SubSkills = subskills
			database.SaveSubskills(skill)
		}
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, skillsCount, skill.Name, exTime)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)

}
