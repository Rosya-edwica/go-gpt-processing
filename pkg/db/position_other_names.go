// position_other_names - получение профессий без других наименований и замена старых наименований на новые
package db

import (
	"fmt"
	"go-gpt-processing/internal/models"
	"go-gpt-processing/pkg/logger"
)

const (
	QuerySelectPositionsWithoutOtherNames = `
		SELECT 
			position.id, position.name
		FROM 
			position as position 
		LEFT JOIN 
			position_to_position as position_to_position 
		ON 
			position_to_position.position_id = position.id 
		WHERE 
			(level = 0 OR level IS NULL)
			AND other_names IS NULL 
		ORDER BY 
			position.id 
		ASC
	`
	QueryUpdatePositionsOtherNames = "UPDATE position SET other_names = '%s' WHERE id=%d"
)

func (d *Database) GetPositionWithoutOtherNames() (positions []models.Position) {
	return d.GetPositionsByQuery(QuerySelectPositionsWithoutOtherNames)
}

func (d *Database) UpdatePositionOtherNames(pos models.Position) {
	query := fmt.Sprintf(QueryUpdatePositionsOtherNames, convertArrayToSQLString(pos.OtherNames), pos.Id)
	d.ExecuteQuery(query)
	logger.LogInfo.Printf("Вариации написания для профессии - %s:%s", pos.Name, convertArrayToSQLString(pos.OtherNames))
}
