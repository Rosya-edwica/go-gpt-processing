package db

import (
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

func (d *Database) GetParentIdForLevelsWithoutExperienceAndSalary() (id int) {
	err := d.Connection.QueryRow(`
		SELECT DISTINCT pos2.id
		FROM test_gpt_position AS pos
		LEFT JOIN test_gpt_position_to_position AS pos_to_pos on pos_to_pos.position_id=pos.id
		LEFT JOIN test_gpt_position as pos2 ON pos2.id = pos_to_pos.parent_position_id
		WHERE pos_to_pos.experience IS NULL 
		AND pos_to_pos.salary IS NULL 
		AND pos_to_pos.position_id IS NOT NULL 
		AND pos_to_pos.level != 0
		LIMIT 1;
	`).Scan(&id)
	checkErr(err)
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
