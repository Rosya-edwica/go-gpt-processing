// Убираем лишние навыки и добавляем новые, если изначально их было мало
package positionsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"strings"
)

type Skills struct {
	Old   []string
	New   []string
	Extra []string
}

func GetExtraSkillsForPosition(positionName string, oldSkills []string) (result Skills, err error) {
	question := fmt.Sprintf(`
	Вырежи все лишние навыки из этого списка "%s" для  профессии "%s" ? В качестве ответа пришли просто список без лишних слов. Если количество нужных навыков будет меньше 10, то добавь от себя правильные hard-skills для этой профессии. Сократи число навыков до 20
	Если список навыков пустой, то предоставь свой
	Предоставь ответ в такой форме, чтобы я мог автоматически его распарсить
	1. Навык 1
	2. Навык 2
	3. Навык 3`, strings.Join(oldSkills, ";"), positionName)
	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return Skills{}, resp.Error
	}
	items := reLines.FindAllString(resp.Answer, -1)
	if len(items) <= 1 {
		return Skills{}, gpt.WrongAnswerError
	}
	var newSkills []string
	for _, i := range items {
		fmt.Println(i, "Sadsadasfasfasf")
		skill := rePointDigit.ReplaceAllString(i, "")
		skill = strings.ReplaceAll(skill, ".", "")
		// Если навык существовал в исходном списке, то сохраняем его отдельно. Тк нужно будет ему поставить is_chatgpt = true
		if checkItemInList(skill, oldSkills) {
			result.Old = append(result.Old, skill)
			// Если навыка не было в исходном списке, значит он новый, поэтому будем строить новую связь профессии и навыка
		} else {
			result.New = append(result.New, skill)
		}
		newSkills = append(newSkills, skill)
	}
	result.Extra = getUniqueSkillsInTwoList(oldSkills, newSkills)

	return result, nil
}

func checkItemInList(item string, list []string) bool {
	for _, i := range list {
		if strings.ToLower(i) == strings.ToLower(item) {
			return true
		}
	}
	return false
}

func getUniqueSkillsInTwoList(oldList []string, newList []string) (list []string) {
	// Выбираем навыки из старого списка, которых нет в новом
	for _, old := range oldList {
		var exist bool
		for _, new := range newList {
			if strings.ToLower(old) == strings.ToLower(new) {
				exist = true
			}
		}
		if !exist {
			list = append(list, old)
		}
	}
	return list
}
