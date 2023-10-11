package db

import (
	"errors"
	"fmt"
	"go-gpt-processing/pkg/logger"
	"go-gpt-processing/pkg/models"
	"strings"
)

func (d *Database) GetPositionWithoutDescription() (positions []models.Position) {
	query := "SELECT id, name FROM test_gpt_position WHERE description IS NULL"
	return d.GetPositionsByQuery(query)
}

func (d *Database) GetPositionWithoutAbout() (positions []models.Position) {
	query := "SELECT id, name FROM test_gpt_position WHERE about IS NULL"
	return d.GetPositionsByQuery(query)
}

func (d *Database) GetPositionWithoutWorkPlaces() (positions []models.Position) {
	query := "SELECT id, name FROM test_gpt_position WHERE work_places IS NULL"
	return d.GetPositionsByQuery(query)
}

func (d *Database) GetPositionWithoutSkills() (positions []models.Position) {
	query := "SELECT id, name FROM test_gpt_position WHERE skills IS NULL"
	return d.GetPositionsByQuery(query)
}

// Получаем список профессий, у которых нет связи с таблицей функций
func (d *Database) GetPositionWithoutFuctions() (positions []models.Position) {
	query := `
		SELECT id, name FROM test_gpt_position
		WHERE id NOT IN (
			SELECT DISTINCT position_id
			FROM test_gpt_position_to_responsibility
		)
	`
	return d.GetPositionsByQuery(query)
}
func (d *Database) GetPositionWithoutOtherNames() (positions []models.Position) {
	query := `
	SELECT 
		position.id, position.name
	FROM 
		test_gpt_position as position 
	LEFT JOIN 
		test_gpt_position_to_position as position_to_position 
	ON 
		position_to_position.position_id = position.id 
	WHERE 
		(level = 0 OR level IS NULL)
		AND other_names IS NULL 
	ORDER BY 
		position.id 
	ASC
	`
	return d.GetPositionsByQuery(query)
}

func (d *Database) GetPositionsByQuery(query string) (positions []models.Position) {
	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		positions = append(positions, models.Position{
			Id:   id,
			Name: name,
		})

	}
	return
}

func (d *Database) GetOnePositionByQuery(query string) (position models.Position) {
	positions := d.GetPositionsByQuery(query)
	if len(positions) == 0 {
		return models.Position{}
	} else {
		return positions[0]
	}
}

func (d *Database) GetOnePositionWithoutEducation() (position models.Position) {
	query := `
		SELECT id, name 
		FROM test_gpt_position 
		WHERE education IS NULL
		LIMIT 1;`
	return d.GetOnePositionByQuery(query)
}

func (d *Database) GetPositionsWithoutLevels() (position []models.Position) {
	query := `
		SELECT 
		pos.id, pos.name
		FROM 
		test_gpt_position as pos
		LEFT JOIN 
		test_gpt_position_to_position as pos_to_pos 
		ON 
		pos_to_pos.position_id = pos.id 
		WHERE 
		(pos_to_pos.level IS NULL OR pos_to_pos.level=0)
		AND pos.education NOT IN ('среднее профессиональное образование', 'без образования| среднее профессиональное образование')
		HAVING 
		(SELECT COUNT(*) FROM test_gpt_position_to_position WHERE parent_position_id=pos.id)=0
		OR 
		(SELECT MAX(level) FROM test_gpt_position_to_position WHERE parent_position_id=pos.id)=0`
	return d.GetPositionsByQuery(query)
}

func (d *Database) GetPositionsLevelsWithoutExperienceAndSalary() (position []models.Position) {
	query := `
		SELECT pos.id, pos.name 
		FROM test_gpt_position AS pos
		LEFT JOIN test_gpt_position_to_position AS pos_to_pos on pos_to_pos.position_id=pos.id
		WHERE pos_to_pos.experience IS NULL 
		AND pos_to_pos.salary IS NULL 
		AND pos_to_pos.position_id IS NOT NULL 
		AND pos_to_pos.level != 0	
		`
	return d.GetPositionsByQuery(query)
}

func (d *Database) GetParentIdsForLevelsWithoutExperienceAndSalary() (ids []int) {
	query := `
		SELECT DISTINCT pos2.id, pos2.name
		FROM test_gpt_position AS pos
		LEFT JOIN test_gpt_position_to_position AS pos_to_pos on pos_to_pos.position_id=pos.id
		LEFT JOIN test_gpt_position as pos2 ON pos2.id = pos_to_pos.parent_position_id
		WHERE pos_to_pos.experience IS NULL 
		AND pos_to_pos.salary IS NULL 
		AND pos_to_pos.position_id IS NOT NULL 
		AND pos_to_pos.level != 0
	`
	positions := d.GetPositionsByQuery(query)
	for _, i := range positions {
		ids = append(ids, i.Id)
	}
	return
}

func (d *Database) GetPositionsLevelsWithoutExperienceAndSalaryByParentId(id int) (positions []models.Position) {
	query := fmt.Sprintf(`
		SELECT pos.id, pos.name 
		FROM test_gpt_position AS pos
		LEFT JOIN test_gpt_position_to_position AS pos_to_pos on pos_to_pos.position_id=pos.id
		WHERE pos_to_pos.experience IS NULL 
		AND pos_to_pos.salary IS NULL 
		AND pos_to_pos.position_id IS NOT NULL 
		AND pos_to_pos.level != 0	
		AND pos_to_pos.parent_position_id = %d
		ORDER BY pos_to_pos.level ASC
		`, id)
	return d.GetPositionsByQuery(query)
}

func (d *Database) CountPositionsWithoutEducation() (count int64) {
	err := d.Connection.QueryRow("SELECT COUNT(*) FROM test_gpt_position WHERE education IS NULL").Scan(&count)
	checkErr(err)
	return
}

func (d *Database) UpdatePositionDescription(pos models.Position) {
	query := fmt.Sprintf(`UPDATE test_gpt_position SET description = '%s' WHERE id=%d`, strings.ReplaceAll(pos.Description, "'", "`"), pos.Id)
	d.ExecuteQuery(query)
	logger.Log.Printf("Полное описание для профессии - %s:%s", pos.Name, pos.Description)
}

func (d *Database) UpdatePositionAbout(pos models.Position) {
	query := fmt.Sprintf(`UPDATE test_gpt_position SET about = '%s' WHERE id=%d`, strings.ReplaceAll(pos.About, "'", "`"), pos.Id)
	d.ExecuteQuery(query)
	logger.Log.Printf("Короткое описание для профессии - %s:%s", pos.Name, pos.About)
}

func (d *Database) UpdatePositionWorkPlaces(pos models.Position) {
	query := fmt.Sprintf(`UPDATE test_gpt_position SET work_places = '%s' WHERE id=%d`, convertArrayToSQLString(pos.WorkPlaces), pos.Id)
	d.ExecuteQuery(query)
	logger.Log.Printf("Места работы для профессии - %s:%s", pos.Name, convertArrayToSQLString(pos.WorkPlaces))
}

func (d *Database) UpdatePositionSkills(pos models.Position) {
	query := fmt.Sprintf(`UPDATE test_gpt_position SET skills = '%s' WHERE id=%d`, convertArrayToSQLString(pos.Skills), pos.Id)
	d.ExecuteQuery(query)
	logger.Log.Printf("Навыки для профессии - %s:%s", pos.Name, convertArrayToSQLString(pos.Skills))

}

func (d *Database) UpdatePositionOtherNames(pos models.Position) {
	query := fmt.Sprintf(`UPDATE test_gpt_position SET other_names = '%s' WHERE id=%d`, convertArrayToSQLString(pos.OtherNames), pos.Id)
	d.ExecuteQuery(query)
	logger.Log.Printf("Вариации написания для профессии - %s:%s", pos.Name, convertArrayToSQLString(pos.OtherNames))
}

// Сохраняем сначала в отдельную таблицу функции
func (d *Database) InsertPositionFunctions(pos models.Position) {
	var functionsIds []int64
	for _, item := range pos.Functions {
		insertQuery := fmt.Sprintf(`INSERT INTO test_gpt_responsibility(name) VALUES('%s')`, strings.ReplaceAll(item, "'", "`"))
		res, err := d.Connection.Exec(insertQuery)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		functionsIds = append(functionsIds, id)
	}
	d.ConnectPositionWithFunctions(pos.Id, functionsIds)
	logger.Log.Printf("Функции для профессии - %s:%s", pos.Name, convertArrayToSQLString(pos.Functions))
}

// После того как сохранили новые функции, проставляем связь с функцией и профессией
func (d *Database) ConnectPositionWithFunctions(posId int, functionsIds []int64) {
	var inserts []string
	for _, function := range functionsIds {
		insertQuery := fmt.Sprintf(`(%d, %d)`, posId, function)
		inserts = append(inserts, insertQuery)
	}
	query := "INSERT INTO test_gpt_position_to_responsibility(position_id, responsibility_id) VALUES " + strings.Join(inserts, ",")
	d.ExecuteQuery(query)
}

func (d *Database) UpdatePositionEducation(pos models.Position) {
	query := fmt.Sprintf(`UPDATE test_gpt_position SET education = '%s' WHERE id=%d`, convertArrayToSQLString(pos.Education), pos.Id)
	d.ExecuteQuery(query)
}

func (d *Database) InsertPositionLevels(position models.Position) {
	for index, item := range position.Levels {
		if len(item.Level) == 0 {
			continue
		}
		insertQuery := fmt.Sprintf("INSERT INTO test_gpt_position(name) VALUES('%s')", strings.ReplaceAll(item.Level, "'", "`"))
		res, err := d.Connection.Exec(insertQuery)
		checkErr(err)
		levelId, err := res.LastInsertId()
		checkErr(err)

		insertQuery = fmt.Sprintf(
			"INSERT INTO test_gpt_position_to_position(position_id, parent_position_id, level, experience, salary) VALUES(%d, %d, %d, '%s', %d)",
			levelId, position.Id, index+1, item.Experience, item.Salary,
		)
		d.ExecuteQuery(insertQuery)

	}
}

func (d *Database) UpdatePositionsLevelExperienceAndSalary(positions []models.Position, parentId int) {
	for _, pos := range positions {
		d.UpdateOnePositionLevelExperienceAndSalary(pos, parentId)
	}
}

func (d *Database) UpdateOnePositionLevelExperienceAndSalary(pos models.Position, parentId int) {
	query := fmt.Sprintf(`
		UPDATE test_gpt_position_to_position
		SET experience = '%s',
		salary = %d
		WHERE position_id = %d
		AND parent_position_id = %d
	`, pos.Experience, pos.Salary, pos.Id, parentId)
	d.ExecuteQuery(query)
}

func (d *Database) GetPositionsByProfArea(area string) (positions []models.Position) {
	query := fmt.Sprintf(`
	SELECT  position.id, position.name from position_to_position 
	left join position on position.id = position_to_position.position_id
	where parent_position_id in
	(SELECT position.id
			FROM position
			LEFT JOIN position_to_prof_area ON position_to_prof_area.position_id=position.id
			LEFT JOIN prof_area_to_specialty ON prof_area_to_specialty.id=position_to_prof_area.area_id
			LEFT JOIN professional_area ON professional_area.id=prof_area_to_specialty.prof_area_id
			WHERE LOWER(professional_area.name) IN ('%s'))

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
