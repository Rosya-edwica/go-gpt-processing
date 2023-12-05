// position_about - получение профессий без короткого описания и замена старого короткого описания на новое
package db

import (
	"fmt"
	"go-gpt-processing/internal/models"
	"go-gpt-processing/pkg/logger"
	"strings"
)

const (
	QuerySelectPositionsWithoutAboutValue = "SELECT id, name FROM position WHERE about IS NULL OR LENGTH(about) = 0"
	QueryUpdatePositionAbout              = "UPDATE position SET about = '%s' WHERE id=%d"
)

func (d *Database) GetPositionWithoutAbout() (positions []models.Position) {
	return d.GetPositionsByQuery(QuerySelectPositionsWithoutAboutValue)
}

func (d *Database) UpdatePositionAbout(pos models.Position) {
	query := fmt.Sprintf(QueryUpdatePositionAbout, strings.ReplaceAll(pos.About, "'", "`"), pos.Id)
	d.ExecuteQuery(query)
	logger.LogInfo.Printf("Короткое описание для профессии - %s:%s", pos.Name, pos.About)
}
