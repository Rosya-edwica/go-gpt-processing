package db

import (
	"fmt"
	"go-gpt-processing/internal/models"
	"go-gpt-processing/pkg/logger"
	"strings"
)

func (d *Database) GetPositionsWithoutDisability() (positions []models.Position) {
	query := `SELECT 
		pos.id, pos.name
		FROM 
		position as pos
		LEFT JOIN 
		position_to_position as pos_to_pos 
		ON 
		pos_to_pos.position_id = pos.id 
		WHERE 
		(pos_to_pos.level IS NULL OR pos_to_pos.level=0)
		AND pos.id NOT IN (
			SELECT DISTINCT position_id FROM position_to_disability
		)`
	return d.GetPositionsByQuery(query)
}

func (d *Database) UpdatePositionDisability(pos models.Position) {
	if len(pos.Disability) == 0 {
		return
	}
	var inserts []string
	for _, i := range pos.Disability {
		insertQuery := fmt.Sprintf(`(%d, %d)`, pos.Id, i.Id)
		inserts = append(inserts, insertQuery)
	}

	query := "INSERT INTO position_to_disability(position_id, disability_id) VALUES " + strings.Join(inserts, ",")
	d.ExecuteQuery(query)
	logger.LogInfo.Printf("Инвалидность для профессии - %s:%#v", pos.Name, pos.Disability)

}

func (d *Database) GetDisablities() (items []models.Disability) {
	rows, err := d.Connection.Query("SELECT id, name FROM disability")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		items = append(items, models.Disability{
			Id:   id,
			Name: name,
		})

	}
	return
}
