package skillsGPT

import (
	"fmt"
	"gpt-skills/gpt"
	"gpt-skills/logger"
	"gpt-skills/models"
	"strings"
)

func CheckSkillIsSoftOrHard(softOrHard string, skill *models.Skill) (err error) {
	var answer, question string
	if softOrHard == "soft" {
		question = fmt.Sprintf(`Look at examples of hard and soft skils. 
		Here are examples of soft skills: Emotional intelligence, Presentation skills, Stress tolerance, Customer centeredness, Resolve disputes and conflicts, Work hard and concentrate even without inspiration, Ability to learn and develop, Keep up to date, Work under tight deadlines, Ability to work hard and concentrate, Speak clear language, Work with information.
		
		Examples of Hard Skills: Excel, Golang, Linux, Legal, Use legal reference systems and court databases, Use accounting software, Document control, Optimize code, Version control, Create infographics, Prepare layout for print, Market and competitor analysis, Specific narrow skills related to the product being sold, Teaching techniques, Develop an educational plan, C++, Git.
		
		Now answer the question Yes or No: Is "%s" a soft-skill?
		`, skill.Name)
	} else {
		question = fmt.Sprintf(`Look at examples of hard and soft skils. 
		Here are examples of soft skills: Emotional intelligence, Presentation skills, Stress tolerance, Customer centeredness, Resolve disputes and conflicts, Work hard and concentrate even without inspiration, Ability to learn and develop, Keep up to date, Work under tight deadlines, Ability to work hard and concentrate, Speak clear language, Work with information.
		
		Examples of Hard Skills: Excel, Golang, Linux, Legal, Use legal reference systems and court databases, Use accounting software, Document control, Optimize code, Version control, Create infographics, Prepare layout for print, Market and competitor analysis, Specific narrow skills related to the product being sold, Teaching techniques, Develop an educational plan, C++, Git.
		
		Now answer the question Yes or No: Is "%s" a hard-skill?`, skill.Name)
	}

	answer, err = gpt.SendRequestToGPT(question)
	if err != nil {
		for {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				return
			}
			answer, err = (gpt.SendRequestToGPT(question))
			if err == nil {
				break
			}
		}
	}
	skill.IsValid, err = gpt.ConvertAnswerToBoolean(answer)
	logger.Log.Printf("Ответ '%s' для вопроса: %s", answer, question)
	return
}
