package skillsGPT

import (
	"fmt"
	"go-gpt-processing/pkg/gpt"
)

func CheckSkillIsSoftOrHard(softOrHard string, name string) (result bool, err error) {
	var answer, question string
	if softOrHard == "soft" {
		question = fmt.Sprintf(`Look at examples of hard and soft skils. 
		Here are examples of soft skills: Emotional intelligence, Presentation skills, Stress tolerance, Customer centeredness, Resolve disputes and conflicts, Work hard and concentrate even without inspiration, Ability to learn and develop, Keep up to date, Work under tight deadlines, Ability to work hard and concentrate, Speak clear language, Work with information.
		
		Examples of Hard Skills: Excel, Golang, Linux, Legal, Use legal reference systems and court databases, Use accounting software, Document control, Optimize code, Version control, Create infographics, Prepare layout for print, Market and competitor analysis, Specific narrow skills related to the product being sold, Teaching techniques, Develop an educational plan, C++, Git.
		
		Now answer the question Yes or No: Is "%s" a soft-skill?
		`, name)
	} else {
		question = fmt.Sprintf(`Look at examples of hard and soft skils. 
		Here are examples of soft skills: Emotional intelligence, Presentation skills, Stress tolerance, Customer centeredness, Resolve disputes and conflicts, Work hard and concentrate even without inspiration, Ability to learn and develop, Keep up to date, Work under tight deadlines, Ability to work hard and concentrate, Speak clear language, Work with information.
		
		Examples of Hard Skills: Excel, Golang, Linux, Legal, Use legal reference systems and court databases, Use accounting software, Document control, Optimize code, Version control, Create infographics, Prepare layout for print, Market and competitor analysis, Specific narrow skills related to the product being sold, Teaching techniques, Develop an educational plan, C++, Git.
		
		Now answer the question Yes or No: Is "%s" a hard-skill?`, name)
	}

	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return false, resp.Error
	}
	result, err = gpt.ConvertAnswerToBoolean(answer)
	return
}
