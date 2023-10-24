// position_functions - получение профессий без функций, сохранение новых функций и связывание профессии с функциями
package db

import (
	"fmt"
	"go-gpt-processing/pkg/logger"
	"go-gpt-processing/pkg/models"
	"strings"
)

const (
	QuerySelectPositionsWithoutFunctions = `
		SELECT id, name FROM position
		WHERE id NOT IN (
			SELECT DISTINCT position_id
			FROM position_to_responsibility
		)
	`
	QueryInsertFunctions                       = "INSERT INTO responsibility(name) VALUES('%s')"
	QueryInsertLinkBetweenPositionAndFunctions = "INSERT INTO position_to_responsibility(position_id, responsibility_id) VALUES "
)

func (d *Database) GetPositionWithoutFuctions() (positions []models.Position) {
	return d.GetPositionsByQuery(QuerySelectPositionsWithoutFunctions)
}

// Сперва нужно сохранить функции
func (d *Database) InsertPositionFunctions(pos models.Position) {
	var functionsIds []int64
	for _, item := range pos.Functions {
		insertQuery := fmt.Sprintf(QueryInsertFunctions, strings.ReplaceAll(item, "'", "`"))
		res, err := d.Connection.Exec(insertQuery)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		functionsIds = append(functionsIds, id)
	}
	d.ConnectPositionWithFunctions(pos.Id, functionsIds)
	logger.LogInfo.Printf("Функции для профессии - %s:%s", pos.Name, convertArrayToSQLString(pos.Functions))
}

// После того как сохранили новые функции, проставляем связь с функцией и профессией
func (d *Database) ConnectPositionWithFunctions(posId int, functionsIds []int64) {
	var inserts []string
	for _, function := range functionsIds {
		insertQuery := fmt.Sprintf(`(%d, %d)`, posId, function)
		inserts = append(inserts, insertQuery)
	}
	query := QueryInsertLinkBetweenPositionAndFunctions + strings.Join(inserts, ",")
	d.ExecuteQuery(query)
}
