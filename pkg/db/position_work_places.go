// position_work_places - получение профессий без мест работы и замена старого места работы на новое
package db

import (
	"fmt"
	"go-gpt-processing/internal/models"
	"go-gpt-processing/pkg/logger"
)

const (
	QuerySelectPositionsWithoutWorkPlaces = "SELECT id, name FROM position WHERE (work_places IS NULL OR LENGTH(work_places) = 0)"
	QueryUpdatePositionWorkPlaces         = "UPDATE position SET work_places = '%s' WHERE id=%d"
)

func (d *Database) GetPositionWithoutWorkPlaces() (positions []models.Position) {
	return d.GetPositionsByQuery(QuerySelectPositionsWithoutWorkPlaces)
}

func (d *Database) UpdatePositionWorkPlaces(pos models.Position) {
	query := fmt.Sprintf(QueryUpdatePositionWorkPlaces, convertArrayToSQLString(pos.WorkPlaces), pos.Id)
	d.ExecuteQuery(query)
	logger.LogInfo.Printf("Места работы для профессии - %s:%s", pos.Name, convertArrayToSQLString(pos.WorkPlaces))
}
