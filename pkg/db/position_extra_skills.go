package db

import (
	"fmt"
	"go-gpt-processing/pkg/gpt/positionsGPT"
	"go-gpt-processing/pkg/models"
	"strings"
)

const (
	QuerySelectPositionsWithoutGPTSkills = `
		SELECT pos.id, pos.name, 
		GROUP_CONCAT(DISTINCT demand.name SEPARATOR '|') as skills
		FROM position AS pos
		LEFT JOIN position_to_demand AS pos_dem ON pos_dem.position_id = pos.id
		LEFT JOIN demand ON demand.id = pos_dem.demand_id
		WHERE pos_dem.is_chatgpt IS NOT TRUE
		GROUP BY pos.id
		ORDER BY skills ASC
	`
)

func (d *Database) GetPositionsWithoutGPTSkills() (positions []models.Position) {
	rows, err := d.Connection.Query(QuerySelectPositionsWithoutGPTSkills)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, skills string

		err = rows.Scan(&id, &name, &skills)
		skills = strings.ReplaceAll(skills, "|", ";")
		positions = append(positions, models.Position{
			Id:     id,
			Name:   name,
			Skills: strings.Split(skills, ";"),
		})
	}
	return
}

func (d *Database) SaveExtraSkills(pos models.Position, skills positionsGPT.Skills) {
	pos.Skills = skills.New
	d.SavePositionSkills(pos)                            // Сохраняем новые навыки
	d.SetPositionSkillsAsGPT(pos.Id, skills.Old)         // Помечаем старые связи как is_chatgpt
	d.RemovePositionLinkWithSkills(pos.Id, skills.Extra) // Удаляем связь профессии с навыками, которые не включил в список GPT (посчитал лишними)
}

func (d *Database) SetPositionSkillsAsGPT(posId int, skills []string) {
	for _, i := range skills {
		query := fmt.Sprintf(`
		UPDATE position_to_demand
		LEFT JOIN demand ON demand.id = position_to_demand.demand_id
		SET is_chatgpt = TRUE
		WHERE LOWER(demand.name) = '%s' and position_id = '%d'`, i, posId)

		d.ExecuteQuery(query)
	}
}

func (d *Database) RemovePositionLinkWithSkills(posId int, skills []string) {
	for _, i := range skills {
		query := fmt.Sprintf(`
		DELETE position_to_demand FROM position_to_demand
		LEFT JOIN demand on demand.id = position_to_demand.demand_id
		WHERE LOWER(demand.name) = '%s' AND position_id = %d`, i, posId)
		d.ExecuteQuery(query)
	}
}
