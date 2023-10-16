// position_experience - получение профессий без опыта или зарплаты и замена старого опыта и образования
package db

import (
	"fmt"
	"go-gpt-processing/pkg/models"
)

const (
	QuerySelectParentIdsWithoutExperienceORSalary = `	
		SELECT DISTINCT pos2.id, pos2.name
		FROM position AS pos
		LEFT JOIN position_to_position AS pos_to_pos on pos_to_pos.position_id=pos.id
		LEFT JOIN position as pos2 ON pos2.id = pos_to_pos.parent_position_id
		WHERE pos_to_pos.experience IS NULL 
		AND pos_to_pos.salary IS NULL 
		AND pos_to_pos.position_id IS NOT NULL 
		AND pos_to_pos.level != 0
	`
	QuerySelectPositionsWithoutExperienceORSalaryByParentId = `
		SELECT pos.id, pos.name 
		FROM position AS pos
		LEFT JOIN position_to_position AS pos_to_pos on pos_to_pos.position_id=pos.id
		WHERE pos_to_pos.experience IS NULL 
		AND pos_to_pos.salary IS NULL 
		AND pos_to_pos.position_id IS NOT NULL 
		AND pos_to_pos.level != 0	
		AND pos_to_pos.parent_position_id = %d
		ORDER BY pos_to_pos.level ASC
	`
	QueryUpdateOnePositionLevelExperienceAndSalary = `
		UPDATE position_to_position
		SET experience = '%s',
		salary = %d
		WHERE position_id = %d
		AND parent_position_id = %d
	`
)

// Берем id родителей, чьи дочерние профессии не имеют опыта или зарплаты
func (d *Database) GetParentIdsForLevelsWithoutExperienceAndSalary() (ids []int) {
	positions := d.GetPositionsByQuery(QuerySelectParentIdsWithoutExperienceORSalary)
	for _, i := range positions {
		ids = append(ids, i.Id)
	}
	return
}

// Берем список дочерних профессий по id родителя
func (d *Database) GetPositionsLevelsWithoutExperienceAndSalaryByParentId(id int) (positions []models.Position) {
	query := fmt.Sprintf(QuerySelectPositionsWithoutExperienceORSalaryByParentId, id)
	return d.GetPositionsByQuery(query)
}

func (d *Database) UpdatePositionsLevelExperienceAndSalary(positions []models.Position, parentId int) {
	for _, pos := range positions {
		d.UpdateOnePositionLevelExperienceAndSalary(pos, parentId)
	}
}

func (d *Database) UpdateOnePositionLevelExperienceAndSalary(pos models.Position, parentId int) {
	query := fmt.Sprintf(QueryUpdateOnePositionLevelExperienceAndSalary, pos.Experience, pos.Salary, pos.Id, parentId)
	d.ExecuteQuery(query)
}
