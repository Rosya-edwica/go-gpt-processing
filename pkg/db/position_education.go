// position_education - получение профессий без образования и замена старого образования на новое
package db

import (
	"fmt"
	"go-gpt-processing/pkg/models"
)

const (
	QuerySelectPositionsWithoutEducationValue      = "SELECT id, name FROM position WHERE education IS NULL OR LENGTH(education) = 0"
	QuerySelectCountPositionsWithoutEducationValue = "SELECT COUNT(*) FROM position WHERE education IS NULL OR LENGTH(education) = 0"
	QueryUpdatePositionEducation                   = "UPDATE position SET education = '%s' WHERE id=%d"
)

func (d *Database) GetOnePositionWithoutEducation() (position models.Position) {
	return d.GetOnePositionByQuery(QuerySelectPositionsWithoutEducationValue)
}

func (d *Database) CountPositionsWithoutEducation() (count int64) {
	err := d.Connection.QueryRow(QuerySelectCountPositionsWithoutEducationValue).Scan(&count)
	checkErr(err)
	return
}

func (d *Database) UpdatePositionEducation(pos models.Position) {
	query := fmt.Sprintf(QueryUpdatePositionEducation, convertArrayToSQLString(pos.Education), pos.Id)
	d.ExecuteQuery(query)
}
