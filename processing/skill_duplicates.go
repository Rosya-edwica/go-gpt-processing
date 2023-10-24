package processing

import (
	"go-gpt-processing/pkg/db"
	"go-gpt-processing/pkg/gpt/skillsGPT"
	"go-gpt-processing/pkg/telegram"
)

func CheckAllSkillsForDuplicates(database *db.Database) {
	const SuccessMessage = "Обработали дубликаты всех навыков"

	pairList := database.GetAllSkillsPair()
	for _, pair := range pairList {
		isDuplicate, err := skillsGPT.CheckSkillsForDuplicates(pair)
		checkErr(err)
		pair.IsDuplicate = isDuplicate
		database.UpdatePair(pair)
	}

	telegram.SuccessMessageMailing(SuccessMessage)
}
