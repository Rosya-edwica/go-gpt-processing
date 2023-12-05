package skill

import (
	"fmt"
	gpt "go-gpt-processing/internal/gpt/skill"
	resp "go-gpt-processing/internal/repositories/skill"
	"go-gpt-processing/pkg/telegram"
	"go-gpt-processing/tools"
)

func CheckAllSkillsForDuplicates(r *resp.Repository) error {
	const SuccessMessage = "Обработали дубликаты всех навыков"
	skills, err := r.GetDuplicates()
	if err != nil {
		return err
	}
	for _, skill := range skills {
		isDuplicate, err := gpt.CheckDuplicates(skill.Name, skill.DuplicateName)
		tools.CheckErr(err)
		skill.IsDuplicate = isDuplicate
		updated, err := r.UpdateDuplicate(skill)
		if err != nil {
			return err
		}
		if !updated {
			fmt.Println("Не удалось обновить дубликаты")
		}
	}
	telegram.SuccessMessageMailing(SuccessMessage)
	return nil
}

func FindTypeGroupForAll(r *resp.Repository) error {
	const SuccessMessage = "Определили тип всех навыков"

	skills, err := r.GetSkillsWithoutTypeGroup()
	if err != nil {
		return err
	}
	for _, skill := range skills {
		groupType, err := gpt.CheckSkillsForTypeGroup(skill.Name)
		tools.CheckErr(err)
		skill.GroupType = groupType
		updated, err := r.UpdateTypeGroup(skill)
		if err != nil {
			return err
		}
		if !updated {
			fmt.Println("Не удалось обновить дубликаты")
		}
	}
	telegram.SuccessMessageMailing(SuccessMessage)
	return nil
}
