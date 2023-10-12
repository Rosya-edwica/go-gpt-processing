package db

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/models"
	"strings"
)

func (d *Database) GetPositionsByProfArea(area string) (positions []models.Position) {
	query := fmt.Sprintf(`
	SELECT  position.id, position.name from position_to_position 
	left join position on position.id = position_to_position.position_id
	WHERE position_to_position.level != 0 
	AND parent_position_id in
	(SELECT position.id
			FROM position
			LEFT JOIN position_to_prof_area ON position_to_prof_area.position_id=position.id
			LEFT JOIN prof_area_to_specialty ON prof_area_to_specialty.id=position_to_prof_area.area_id
			LEFT JOIN professional_area ON professional_area.id=prof_area_to_specialty.prof_area_id
			WHERE LOWER(professional_area.name) IN ('%s'))
			AND position.id NOT IN (
				SELECT DISTINCT position_id FROM position_to_demand)

	`, strings.ToLower(area))

	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		positions = append(positions, models.Position{
			Id:       id,
			Name:     name,
			ProfArea: area,
		})

	}
	return
}

func (d *Database) SaveDemand(name string) (demandId int64, err error) {
	demandId = d.CheckDemandExist(name)
	if demandId != 0 {
		return
	}
	query := fmt.Sprintf("INSERT INTO demand(name) VALUES('%s')", strings.ReplaceAll(name, "'", "`"))
	res, err := d.Connection.Exec(query)
	if err != nil {
		return 0, err
	}

	demandId, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if demandId == 0 {
		return 0, err
	}

	return demandId, nil

}

func (d *Database) CheckDemandExist(name string) (demandId int64) {
	lowerName := strings.ToLower(name)
	query := fmt.Sprintf("SELECT id FROM demand WHERE LOWER(name) = '%s'", lowerName)
	err := d.Connection.QueryRow(query).Scan(&demandId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0
		}
		return 0
	}
	return
}

func (d *Database) SavePositionSkills(pos models.Position) (err error) {
	var skillsIds []int64
	for _, skill := range pos.Skills {
		id, err := d.SaveDemand(skill)
		if err != nil {
			fmt.Println("ERROR: ", err)
			continue
		}
		skillsIds = append(skillsIds, id)
	}
	if len(skillsIds) == 0 {
		text := fmt.Sprintf("Не получилось сохранить навыки для профессии:'%s', т.к. количество навыков равно 0", pos.Name)
		return errors.New(text)
	}

	d.JoinPositionWithSkills(pos.Id, skillsIds)
	return nil

}

func (d *Database) JoinPositionWithSkills(posId int, skillsIds []int64) {
	var inserts []string
	for _, id := range skillsIds {
		insertQuery := fmt.Sprintf(`(%d, %d, 1, true, false)`, posId, id)
		inserts = append(inserts, insertQuery)
	}
	query := "INSERT INTO position_to_demand(position_id, demand_id, is_custom, is_chatgpt, is_delete) VALUES " + strings.Join(inserts, ",")
	d.ExecuteQuery(query)
	fmt.Println(posId)
	fmt.Println("Success!")

}

func (d *Database) GetProfAreaList() (areas []string) {
	query := "SELECT name FROM professional_area"
	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		checkErr(err)
		areas = append(areas, name)
	}
	return
}
