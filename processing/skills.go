package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"strings"
)

func CheckAllSkillsForDuplicates(database *db.Database) {
	fmt.Println("Ищем дубликаты...")
	for {
		pair := database.GetSkillsPair()
		if pair.Id == 0 {
			return
		}
		err := skillsGPT.CheckSkillsForDuplicates(&pair)
		checkErr(err)
		if err == nil {
			database.UpdatePair(pair)
		}
	}
}

func CheckAllSkillsForSoftOrHardSkill(database *db.Database, softOrHard string) {
	fmt.Printf("Ищем %s-skills...\n", softOrHard)
	for {
		skill := database.GetSkill(softOrHard)
		if skill.Id == 0 {
			return
		}
		err := skillsGPT.CheckSkillIsSoftOrHard(softOrHard, &skill)
		checkErr(err)
		if err == nil {
			database.UpdateSkill(softOrHard, skill)
		}
		Pause(5)
	}
}

func CheckAllSkillsForGroupType(database *db.Database) {
	fmt.Println("Пытаемся определить группу для навыка")
	for {
		skill := database.GetSkillWithoutGroup()
		if skill.Id == 0 {
			return
		}
		err := skillsGPT.CheckSkillsForTypeGroup(&skill)
		checkErr(err)
		if err == nil {
			database.UpdateSkillGroup(skill)
		}
		Pause(5)
	}
}

func CollectForAllSkillsSubSkills(database *db.Database) {
	fmt.Println("Ищем поднавыки")
	skills := database.GetSkills()
	for i, skill := range skills {
		subskills, err := skillsGPT.GetSubSkills(skill.Name)
		checkErr(err)
		if len(subskills) != 0 {
			skill.SubSkills = subskills
			database.SaveSubskills(skill)
		}
		fmt.Printf("[%d/%d] Поднавыки для навыка - %s:\n %s\n\n", i+1, len(skills), skill.Name, strings.Join(skill.SubSkills, "|"))
		Pause(5)
	}

}

func CollectForAllSkillsTests(database *db.Database) {
	fmt.Println("Ищем тесты для навыков")
	skills := database.GetSkills()
	for i, skill := range skills {
		test, err := skillsGPT.GetTestForSkill(skill.Name)
		checkErr(err)
		database.SaveSkillsTest(skill.Id, test)
		fmt.Printf("[%d/%d] Тесты для навыка - %s:\n", i+1, len(skills), skill.Name)
	}
}
