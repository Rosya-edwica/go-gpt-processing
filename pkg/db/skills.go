package db

import (
	"fmt"
	"go-gpt-processing/pkg/models"
	"strings"
)

func (d *Database) GetSkillsPair() (skills models.Skill) {
	query := `
		SELECT id, demand_name, dup_demand_name
		FROM demand_duplicate
		WHERE is_duplicate_gpt IS NULL
		LIMIT 1`

	rows, err := d.Connection.Query(query)
	checkErr(err)
	for rows.Next() {
		var first, dupName string
		var id int

		err = rows.Scan(&id, &first, &dupName)
		skills = models.Skill{
			Id:            id,
			Name:          first,
			DuplicateName: dupName,
		}
	}
	return
}

func (d *Database) UpdatePair(skills models.Skill) {
	query := fmt.Sprintf(`
		UPDATE demand_duplicate
		SET is_duplicate_gpt = %t
		WHERE id = %d`, skills.IsDuplicate, skills.Id)
	d.ExecuteQuery(query)
}

func (d *Database) GetSkill(softOrHard string) (skill models.Skill) {
	var query string
	if softOrHard == "soft" {
		query = `
		SELECT id, translated
		FROM demand
		WHERE is_soft_gpt IS NULL AND is_hard_gpt IS NOT TRUE AND is_custom IS NOT TRUE AND type_group = 'навык' AND is_deleted IS FALSE
		LIMIT 1`
	} else {
		query = `
		SELECT id, translated
		FROM demand
		WHERE is_soft_gpt IS NOT TRUE AND is_hard_gpt IS NULL AND is_custom IS NOT TRUE AND type_group = 'навык' AND is_deleted IS FALSE
		LIMIT 1`
	}

	rows, err := d.Connection.Query(query)
	checkErr(err)
	for rows.Next() {
		var name string
		var id int

		err = rows.Scan(&id, &name)
		skill = models.Skill{
			Id:   id,
			Name: name,
		}
	}
	return
}

func (d *Database) UpdateSkill(softOrHard string, skill models.Skill) {
	var query string
	if softOrHard == "soft" {
		query = fmt.Sprintf(`
			UPDATE demand
			SET is_soft_gpt = %t
			WHERE id = %d`, skill.IsValid, skill.Id)
	} else {
		query = fmt.Sprintf(`
			UPDATE demand
			SET is_hard_gpt = %t
			WHERE id = %d`, skill.IsValid, skill.Id)
	}
	d.ExecuteQuery(query)

}

func (d *Database) GetSkillWithoutGroup() (skill models.Skill) {
	query := `
		SELECT id, translated
		FROM demand
		WHERE translated IS NOT NULL AND type_group IS NULL
		LIMIT 1`
	rows, err := d.Connection.Query(query)
	checkErr(err)
	for rows.Next() {
		var name string
		var id int

		err = rows.Scan(&id, &name)
		skill = models.Skill{
			Id:   id,
			Name: name,
		}
	}
	return
}

func (d *Database) UpdateSkillGroup(skill models.Skill) {
	query := fmt.Sprintf(`
		UPDATE demand
		SET type_group='%s'
		WHERE lower(translated) = '%s' `, skill.GroupType, strings.ToLower(skill.Name))
	d.ExecuteQuery(query)
}

func (d *Database) GetSkills() (skills []models.Skill) {
	query := `
		SELECT id, name
		FROM demand
		WHERE type_group = 'навык' and is_hard_gpt is true AND is_deleted is false`
	rows, err := d.Connection.Query(query)
	checkErr(err)
	for rows.Next() {
		var name string
		var id int

		err = rows.Scan(&id, &name)
		skills = append(skills, models.Skill{
			Id:   id,
			Name: name,
		})
	}
	return
}

// Передаешь название колонки и функция вернет все строки, где значение этой колонки равно NULL
func (d *Database) GetNullableSkillsInColumn(column string) (skills []models.Skill) {
	query := fmt.Sprintf(`
	SELECT id, name
	FROM demand
	WHERE %s IS NULL
	`, column)
	return d.GetSkillsByQuery(query)

}

func (d *Database) GetSkillsByQuery(query string) (skills []models.Skill) {
	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)

		skill := models.Skill{
			Id:   id,
			Name: name,
		}
		skills = append(skills, skill)
	}
	return
}

func (d *Database) SaveSubskills(skill models.Skill) {
	var subSkillsIds []int64
	for _, item := range skill.SubSkills {
		insertQuery := fmt.Sprintf(`INSERT INTO subskills(name) VALUES(%s)`, strings.ReplaceAll(item, `'`, ``))
		res, err := d.Connection.Exec(insertQuery)
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)
		subSkillsIds = append(subSkillsIds, id)
	}
	d.ConnectSkillsWithSubSkills(skill.Id, subSkillsIds)
	fmt.Sprintln("Success")
}

func (d *Database) ConnectSkillsWithSubSkills(skillId int, subskillsIds []int64) {
	var inserts []string
	for _, id := range subskillsIds {
		insertQuery := fmt.Sprintf(`(%d, %d)`, skillId, id)
		inserts = append(inserts, insertQuery)
	}
	query := "INSERT INTO demand_to_subskills(demand_id, subskill_id) VALUES" + strings.Join(inserts, ",")
	d.ExecuteQuery(query)
	fmt.Sprintln(query)

}

func (d *Database) SaveSkillsTest(skillId int, test models.Test) {
	for _, question := range test.Questions {
		query := fmt.Sprintf(`
			INSERT IGNORE INTO demand_tests(demand_id, question, choices, answer)
			VALUES (%d, '%s', '%s', '%s');`,
			skillId, question.Text, strings.Join(question.Choices, "|"), question.Answer)
		d.ExecuteQuery(query)
	}
}

func (d *Database) UpdateSkillColumn(skillId int, column string, value string) {
	query := fmt.Sprintf(`
	UPDATE demand
	SET %s = '%s'
	WHERE id = %d	
	`, column, value, skillId)

	d.ExecuteQuery(query)
}
