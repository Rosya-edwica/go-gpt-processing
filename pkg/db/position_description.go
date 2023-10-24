// position_description - получение профессий без описания и замена старого описания на новое
package db

import (
	"fmt"
	"go-gpt-processing/pkg/logger"
	"go-gpt-processing/pkg/models"
	"strings"
)

const (
	QuerySelectPositionWithoutDescriptionValue = "SELECT id, name FROM position WHERE description IS NULL OR LENGTH(description) = 0"
	QueryUpdatePositionDescription             = "UPDATE position SET description = '%s' WHERE id=%d"
)

func (d *Database) GetPositionWithoutDescription() (positions []models.Position) {
	return d.GetPositionsByQuery(QuerySelectPositionWithoutDescriptionValue)
}

func (d *Database) UpdatePositionDescription(pos models.Position) {
	query := fmt.Sprintf(QueryUpdatePositionDescription, strings.ReplaceAll(pos.Description, "'", "`"), pos.Id)
	d.ExecuteQuery(query)
	logger.LogInfo.Printf("Полное описание для профессии - %s:%s", pos.Name, pos.Description)

}
