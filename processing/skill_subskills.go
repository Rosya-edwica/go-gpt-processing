package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CollectForAllSkillsSubSkills(database *db.Database) {
	const SuccessMessage = "Подобрали поднавыки для всех навыков"

	skills := database.GetSkills()

	for _, skill := range skills {
		subskills, err := skillsGPT.GetSubSkills(skill.Name)
		checkErr(err)
		if len(subskills) != 0 {
			skill.SubSkills = subskills
			database.SaveSubskills(skill)
		}
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)

}
