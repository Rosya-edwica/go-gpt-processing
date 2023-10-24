package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CollectForAllSkillsTests(database *db.Database) {
	const SuccessMessage = "Подобрали тесты для всех навыков"

	skills := database.GetSkills()
	for _, skill := range skills {
		test, err := skillsGPT.GetTestForSkill(skill.Name)
		checkErr(err)
		database.SaveSkillsTest(skill.Id, test)
	}
	telegram.SuccessMessageMailing(SuccessMessage)
}
