package positionsGPT

import (
	"errors"
	"fmt"
	"gpt-skills/gpt"
	"strings"
)

func GetAboutForPosition(name string) (about string, err error) {
	question := fmt.Sprintf(`Составь описание профессии "%s" в одну строчку`, name)
	about, err = gpt.SendRequestToGPT(question)
	if about == "" {
		return "", errors.New(fmt.Sprintf("Пустое описание для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(about), "я не могу") {
		return "", errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", about, name))
	}
	return
}

func GetDescriptionForPosition(name string) (descr string, err error) {
	question := fmt.Sprintf(`Составь подробное описание профессии "%s"`, name)
	descr, err = gpt.SendRequestToGPT(question)
	if descr == "" {
		return "", errors.New(fmt.Sprintf("Пустое описание для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(descr), "я не могу") {
		return "", errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", descr, name))
	}
	return
}

func GetOtherNamesForPosition(name string) (otherNames []string, err error) {
	question := fmt.Sprintf(`Составь список из 30 вариантов написания профессии "%s". Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	answer, err := gpt.SendRequestToGPT(question)
	otherNames = strings.Split(answer, ",")

	if len(otherNames) <= 1 {
		return nil, errors.New(fmt.Sprintf("Не удалось поделить ответ по запятым: %s", answer))
	}
	if answer == "" {
		return nil, errors.New(fmt.Sprintf("Пустое описание для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", answer, name))
	}
	return
}

func GetWorkPlacesForPosition(name string) (workPlaces []string, err error) {
	question := fmt.Sprintf(`Составь список из 6 мест где может работать %s . Пиши в строчку. Не используй нумерацию. В качестве разделителя используй знак ,`, name)
	answer, err := gpt.SendRequestToGPT(question)
	workPlaces = strings.Split(answer, ",")

	if len(workPlaces) <= 1 {
		return nil, errors.New(fmt.Sprintf("Не удалось поделить ответ по запятым: %s", answer))
	}
	if answer == "" {
		return nil, errors.New(fmt.Sprintf("Пустое ответ для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", answer, name))
	}
	return
}


func GetSkillsForPosition(name string) (skills []string, err error) {
	question := fmt.Sprintf(`
	Составь список из 90 профессиональных навыков и знаний для профессии %s.  Пиши в строчку. Не используй нумерацию. 
	Не используй в описании такие слова как: знание, умение, владение, работа с, понимание. В качестве разделителя используй знак ,`, 
	name)
	answer, err := gpt.SendRequestToGPT(question)
	skills = strings.Split(answer, ",")

	if len(skills) <= 1 {
		return nil, errors.New(fmt.Sprintf("Не удалось поделить ответ по запятым: %s", answer))
	}
	if answer == "" {
		return nil, errors.New(fmt.Sprintf("Пустое ответ для профессии: %s", name))
	} else if strings.Contains(strings.ToLower(answer), "я не могу") {
		return nil, errors.New(fmt.Sprintf("Неправильный ответ '%s' для профессии - %s", answer, name))
	}
	return
}