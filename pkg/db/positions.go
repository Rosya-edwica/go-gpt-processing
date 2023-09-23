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

func (d *Database) GetPositionWithoutFuctions() (positions []models.Position) {
	query := ""
	return d.GetPositionsByQuery(query)
}
func (d *Database) GetPositionWithoutOtherNames() (positions []models.Position) {
	query := `
		SELECT pos.id, pos.name FROM test_gpt_position as pos
		LEFT JOIN test_gpt_position_to_position as pos_to_pos on pos_to_pos.position_id = pos.id
		WHERE pos_to_pos.level = 0 AND other_names IS NULL
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

func (d *Database) UpdatePositionFunctions(pos models.Position) {
	query := ""
	d.ExecuteQuery(query)
	logger.Log.Printf("Функции для профессии - %s:%s", pos.Name, convertArrayToSQLString(pos.Functions))
}

func convertArrayToSQLString(items []string) (result string) {
	result = strings.Join(items, "|")
	result = strings.ReplaceAll(result, ".", "")
	result = strings.ToLower(result)
	result = strings.ReplaceAll(result, "'", "`")
	return
}
