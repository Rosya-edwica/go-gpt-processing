package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CheckAllSkillsForSoftOrHardSkill(database *db.Database, softOrHard string) {
	const SuccessMessage = "Закончили проверку навыков на hard/soft-skills"
	var op = "processing.skill_hard_soft"

	skills := database.GetSkillsHardSkills(softOrHard)
	skillsCount := len(skills)
	for i, skill := range skills {
		isTrue, exTime, err := skillsGPT.CheckSkillIsSoftOrHard(softOrHard, skill.Name)
		checkErr(err)
		skill.IsValid = isTrue
		database.UpdateSkill(softOrHard, skill)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, skillsCount, skill.Name, exTime)
		Pause(5)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
