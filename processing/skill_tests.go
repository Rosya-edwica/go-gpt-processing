package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CollectForAllSkillsTests(database *db.Database) {
	const SuccessMessage = "Подобрали тесты для всех навыков"
	var op = "processing.skill_tests"

	skills := database.GetSkills()
	skillsCount := len(skills)
	for i, skill := range skills {
		test, exTime, err := skillsGPT.GetTestForSkill(skill.Name)
		checkErr(err)
		database.SaveSkillsTest(skill.Id, test)
		fmt.Printf("%s\t[%d/%d] %s (Time: %d s)\n", op, i+1, skillsCount, skill.Name, exTime)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
