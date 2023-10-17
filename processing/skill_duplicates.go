package processing

import (
	"fmt"
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CheckAllSkillsForDuplicates(database *db.Database) {
	const SuccessMessage = "Обработали дубликаты всех навыков"
	var op = "processing.skill_duplicates"

	pairList := database.GetAllSkillsPair()
	pairCount := len(pairList)
	for i, pair := range pairList {
		isDuplicate, exTime, err := skillsGPT.CheckSkillsForDuplicates(pair)
		checkErr(err)
		pair.IsDuplicate = isDuplicate
		database.UpdatePair(pair)
		fmt.Printf("%s\t[%d/%d] %s|%s (Time: %d s)\n", op, i+1, pairCount, pair.Name, pair.DuplicateName, exTime)
	}

	telegram.SuccessMessageMailing(SuccessMessage)
}
