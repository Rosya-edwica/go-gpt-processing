package skillsGPT

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/gpt"
	"go-gpt-processing/pkg/logger"
	"go-gpt-processing/pkg/models"
	"strings"
)

func CheckSkillsForTypeGroup(skill *models.Skill) (err error) {
	question := fmt.Sprintf(` Identify which category the word "%s" belongs to. 
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
	`, skill.Name)
	answer, err := gpt.SendRequestToGPT(strings.TrimSpace(question))
	fmt.Println(answer)
	if err != nil {
		return errors.New(fmt.Sprintf("ОШИБКА: %s", err.Error()))
	}
	if strings.Contains(answer, "профессия/специальность/должность") || strings.Contains(answer, "2") {
		skill.GroupType = "профессия"
	} else if strings.Contains(answer, "навык") || strings.Contains(answer, "1") {
		skill.GroupType = "навык"
	} else if strings.Contains(answer, "другое") || strings.Contains(answer, "3") {
		skill.GroupType = "другое"
	} else {
		return errors.New(fmt.Sprintf("Неправильный ответ: ответ - %s. вопрос: %s", answer, question))
	}

	logger.Log.Printf("Ответ '%s' для вопроса: %s", answer, question)

	return
}
