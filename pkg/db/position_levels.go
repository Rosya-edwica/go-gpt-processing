// position_levels - получение профессий без дочерних профессий и сохранение новых профессий вместе с ЗП и опытом
package db

import (
	"fmt"
	"go-gpt-processing/internal/models"
	"strings"
)

const (
	QuerySelectPositionsWithoutLevels = `
		SELECT 
		pos.id, pos.name
		FROM 
		position as pos
		LEFT JOIN 
		position_to_position as pos_to_pos 
		ON 
		pos_to_pos.position_id = pos.id 
		WHERE 
		(pos_to_pos.level IS NULL OR pos_to_pos.level=0)
		AND
		(pos.education NOT IN ('среднее профессиональное образование', 'без образования| среднее профессиональное образование') OR education IS NULL)
		HAVING 
		(SELECT COUNT(*) FROM position_to_position WHERE parent_position_id=pos.id)=0
		OR 
		(SELECT MAX(level) FROM position_to_position WHERE parent_position_id=pos.id)=0
	`
	QueryInsertNewPositions                 = "INSERT INTO position(name) VALUES('%s')"
	QueryInsertLinkBetweenPositionAndLevels = "INSERT INTO position_to_position(position_id, parent_position_id, level, experience, salary) VALUES(%d, %d, %d, '%s', %d)"
)

func (d *Database) GetPositionsWithoutLevels() (position []models.Position) {
	return d.GetPositionsByQuery(QuerySelectPositionsWithoutLevels)
}

func (d *Database) InsertPositionLevels(position models.Position) {
	for index, item := range position.Levels {
		if len(item.Level) == 0 {
			continue
		}
		insertQuery := fmt.Sprintf(QueryInsertNewPositions, strings.ReplaceAll(item.Level, "'", "`"))
		res, err := d.Connection.Exec(insertQuery)
		checkErr(err)
		levelId, err := res.LastInsertId()
		checkErr(err)

		insertQuery = fmt.Sprintf(QueryInsertLinkBetweenPositionAndLevels, levelId, position.Id, index+1, item.Experience, item.Salary)
		d.ExecuteQuery(insertQuery)

	}
}
