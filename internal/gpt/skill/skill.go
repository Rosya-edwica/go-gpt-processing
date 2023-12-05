package skill

import (
	"fmt"
	"go-gpt-processing/internal/gpt"
	"go-gpt-processing/internal/models"
	"regexp"
	"strings"
)

var (
	reLines              = regexp.MustCompile(`\d+. .*?\n`)
	reDigits             = regexp.MustCompile(`\d+. |\n`)
	testTitleRegexp      = regexp.MustCompile(`\d+. .*?\n`)
	testTitleRegexpSub   = regexp.MustCompile(`\d+.|\n|Вопрос:`)
	testChoicesRegexp    = regexp.MustCompile(`\w\) .*`)
	testChoicesRegexpSub = regexp.MustCompile(`\w\) |\n|\( ответ \)`)
	testAnswerRegexp     = regexp.MustCompile(`Ответ:.*|\w\) .*\( ответ \)`)
	testAnswerRegexpSub  = regexp.MustCompile(`Ответ: |\n|\( ответ \)|\w\) `)
)

func CheckDuplicates(first, second string) (isDuplicate bool, err error) {
	question := fmt.Sprintf("Можно ли считать эти навыки дубликатами: '%s' и '%s'? Ответь Да или Нет.", first, second)
	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return false, resp.Error
	}

	isDuplicate, err = gpt.ConvertAnswerToBoolean(resp.Answer)
	return
}

func CheckSkillsForTypeGroup(name string) (groupType string, err error) {
	question := fmt.Sprintf(`Identify which category the word "%s" belongs to. 
	If it is a skill, ability or competence, mark it as 1. 
	If it is a profession or position, mark as 2. Otherwise, mark as 3.

	Here's an example:
	The skill, ability or competency category includes the following - Python, Golang, GIT, creativity,	yii2 proficiency, 
	administration, legal opinion, knowledge of vehicle design, effective time management,
	Treatment and diagnostic process steps, programming, optimization, team organization, paperwork, 
	motivation, Spanish, Epidural anesthesia, Knowledge of basic categories of pedagogy, Boroscopic examination, 
	Development of proposals to improve the reliability of operating equipment, Verification of personnel compliance with regulations 
	operation of the equipment, Forecasting of optimal well flow rate, Rules of safe organization of labor 
	When manufacturing carpentry products, Determine malfunctions in the operation of monorail carts, Observing the readings of the
	Control and measuring instruments, Selection of seed varieties, Improvement of technical condition of land reclamation systems, 
	Soil cultivation, Bitrix 24, 1C, Autocad, Scrum, MS Excel, MS Access, Access, Excel, 1C: Enterprise 8, 
	Sublime text, Labor Code of the Russian Federation, Diligence, Good learning ability, Car repair, Transportation logistics, 
	SAP, Telemarketing, Data Analysis, SEO-Promotion, Fire Risk Calculation, Mining, Teaching.


	The category "profession or position" includes - Sailor, accountant, programmer, designer, python developer, 
	consultant, Visual Arts teacher, Engineer for mechanization and automation of production processes, 
	Spare Parts Purchasing Manager, Leading Manager of Inbound Tourism, Stain Remover, Receptionist, Accommodation Worker, 
	Automobile repair mechanic 2nd grade, Physical-mechanical testing laboratory technician 5th grade, 
	Warranty Engineer (Regional Supervisor), 5th Grade Mine Furnace Gas Engineer, Dentist, Doctor, Teacher, Chemist, 
	Loader, Welder, Cook, 1C-programmer, HR-manager, Recruiter, Recruitment specialist, Logistician, CEO, Director, 
	Manager, Auditor, Inspector, Insurer, Financier, Economist, Producer, Legal Consultant, Lawyer, Taxi Driver, Policeman, 
	SAP specialist, Teacher.


	The category "Other" includes - driver card, secondary education, category C, category B, Shift work, no experience,
	experience in landscaping, filming, staffing, working on staff, working on wear and tear, curtain, tanker, 
	work without experience, Law, Advertising, Import, IT field, military ID, regional hiring, warranty, maintenance, A4 sheet of paper,
	desire to earn, special equipment, tractor, car, cargo, tractor, neat appearance, Civil Defense, Ministry of Internal Affairs, 
	Desire to work in a combat unit, sports nutrition, experience in active sales, 
	2nd Special Police Regiment of the Main Department of the Ministry of Internal Affairs of Russia in Moscow, speed of work, Gosts, Teaching staff.
	`, name)
	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return "", resp.Error
	}
	if strings.Contains(resp.Answer, "профессия/специальность/должность") || strings.Contains(resp.Answer, "2") {
		groupType = "профессия"
	} else if strings.Contains(resp.Answer, "навык") || strings.Contains(resp.Answer, "1") {
		groupType = "навык"
	} else if strings.Contains(resp.Answer, "другое") || strings.Contains(resp.Answer, "3") {
		groupType = "другое"
	} else {
		return "", gpt.WrongAnswerError
	}
	return
}

func FindSubSkills(query string) (skills []string, err error) {
	question := fmt.Sprintf(`
	какими hard-skills нужно обладать, чтобы изучить "%s" - сократи ответ до перечня навыков
	`, query)

	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return nil, resp.Error
	}
	lines := reLines.FindAllString(resp.Answer+"\n", -1)
	for _, line := range lines {
		skill := reDigits.ReplaceAllString(line, "")
		skills = append(skills, skill)
	}
	return
}

func CheckSoftSkill(name string) (bool, error) {
	question := fmt.Sprintf(`Look at examples of hard and soft skils. 
		Here are examples of soft skills: Emotional intelligence, Presentation skills, Stress tolerance, Customer centeredness, Resolve disputes and conflicts, Work hard and concentrate even without inspiration, Ability to learn and develop, Keep up to date, Work under tight deadlines, Ability to work hard and concentrate, Speak clear language, Work with information.
		
		Examples of Hard Skills: Excel, Golang, Linux, Legal, Use legal reference systems and court databases, Use accounting software, Document control, Optimize code, Version control, Create infographics, Prepare layout for print, Market and competitor analysis, Specific narrow skills related to the product being sold, Teaching techniques, Develop an educational plan, C++, Git.
		
		Now answer the question Yes or No: Is "%s" a soft-skill?
	`, name)
	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return false, resp.Error
	}
	result, err := gpt.ConvertAnswerToBoolean(resp.Answer)
	return result, err
}

func CheckHardSkill(name string) (bool, error) {
	question := fmt.Sprintf(`Look at examples of hard and soft skils. 
		Here are examples of soft skills: Emotional intelligence, Presentation skills, Stress tolerance, Customer centeredness, Resolve disputes and conflicts, Work hard and concentrate even without inspiration, Ability to learn and develop, Keep up to date, Work under tight deadlines, Ability to work hard and concentrate, Speak clear language, Work with information.
		
		Examples of Hard Skills: Excel, Golang, Linux, Legal, Use legal reference systems and court databases, Use accounting software, Document control, Optimize code, Version control, Create infographics, Prepare layout for print, Market and competitor analysis, Specific narrow skills related to the product being sold, Teaching techniques, Develop an educational plan, C++, Git.
		
		Now answer the question Yes or No: Is "%s" a hard-skill?
	`, name)
	resp := gpt.SendRequestToGPT(question)
	if resp.Error != nil {
		return false, resp.Error
	}
	result, err := gpt.ConvertAnswerToBoolean(resp.Answer)
	return result, err
}

func GetTestForSkill(query string) (test models.Test, err error) {
	gptQuestion := fmt.Sprintf("Для изучения навыка '%s' составь тест из 10 вопросов с вариантами ответов в такой структуре: 1. Вопрос\na) первый вариант\nb) второй вариант\nc) третий вариант\nd) четвертый вариант\nОтвет: полный вариант", query)
	resp := gpt.SendRequestToGPT(gptQuestion)
	if resp.Error != nil {
		return models.Test{}, resp.Error
	}
	test = parseTest(strings.TrimSpace(resp.Answer))
	return
}

func parseTest(text string) (test models.Test) {
	listQuestionsText := strings.Split(text, "\n\n")
	for _, item := range listQuestionsText {
		question := parseQuestion(item)
		test.Questions = append(test.Questions, question)
	}
	return
}

func parseQuestion(text string) (question models.Question) {
	question.Text = testTitleRegexp.FindString(text)
	question.Text = testTitleRegexpSub.ReplaceAllString(question.Text, "")
	question.Text = strings.TrimSpace(question.Text)

	question.Choices = parseChoices(text)
	question.Answer = parseAnswer(text)
	return
}

func parseChoices(text string) (choices []string) {
	choicesFinded := testChoicesRegexp.FindAllString(text, -1)
	for _, item := range choicesFinded {
		item = testChoicesRegexpSub.ReplaceAllString(item, "")
		item = strings.TrimSpace(item)
		choices = append(choices, item)
	}
	return
}
func parseAnswer(text string) (answer string) {
	answer = testAnswerRegexp.FindString(text)
	answer = testAnswerRegexpSub.ReplaceAllString(answer, "")
	answer = strings.TrimSpace(answer)
	return
}

func GetDescriptionForSkill(name string) (description string, err error) {
	question := fmt.Sprintf("Опиши этот навык - '%s'", name)
	resp := gpt.SendRequestToGPT(question)
	return resp.Answer, resp.Error
}
