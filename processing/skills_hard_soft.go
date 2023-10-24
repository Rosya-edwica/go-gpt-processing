package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CheckAllSkillsForSoftOrHardSkill(database *db.Database, softOrHard string) {
	const SuccessMessage = "Закончили проверку навыков на hard/soft-skills"

	skills := database.GetSkillsHardSkills(softOrHard)
	for _, skill := range skills {
		isTrue, err := skillsGPT.CheckSkillIsSoftOrHard(softOrHard, skill.Name)
		checkErr(err)
		skill.IsValid = isTrue
		database.UpdateSkill(softOrHard, skill)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
