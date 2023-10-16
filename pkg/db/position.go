package db

import (
	"go-gpt-processing/pkg/models"
)

func (d *Database) GetOnePositionByQuery(query string) (position models.Position) {
	positions := d.GetPositionsByQuery(query)
	if len(positions) == 0 {
		return models.Position{}
	} else {
		return positions[0]
	}
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
