package db

import (
	"fmt"
	"gpt-skills/models"
	"strings"
)

func (d *Database) ExecuteQuery(query string) {
	tx, _ := d.Connection.Begin()
	_, err := d.Connection.Exec(query)
	checkErr(err)
	tx.Commit()
}

func (d *Database) GetSkillsPair() (pair models.Pair) {
	query := `
		SELECT id, demand_name, dup_demand_name, is_duplicate
		FROM demand_duplicate
		WHERE is_duplicate_gpt IS NULL
		LIMIT 1`
	
	rows, err := d.Connection.Query(query)
	checkErr(err)
	for rows.Next() {
		var first, second string
		var id int
		var isDuplicate bool

		err = rows.Scan(&id, &first, &second, &isDuplicate)
		pair = models.Pair{
			First: first,
			Second: second,
			Id: id,
			IsDuplicate: isDuplicate,
		}
	}
	return
}


func (d *Database) UpdatePair(pair models.Pair) {
	query := fmt.Sprintf(`
		UPDATE demand_duplicate
		SET is_duplicate_gpt = %t
		WHERE id = %d`, pair.IsDuplicate, pair.Id)
	d.ExecuteQuery(query)
}

func (d *Database) GetSkill(softOrHard string) (skill models.Skill) {
	var query string
	if softOrHard == "soft" {
		query = `
		SELECT id, name, is_displayed
		FROM demand
		WHERE is_soft_gpt IS NULL AND is_hard_gpt IS NOT TRUE AND is_custom IS NOT TRUE
		LIMIT 1`
	} else {
		query = `
		SELECT id, name, is_displayed
		FROM demand
		WHERE is_soft_gpt IS NOT TRUE AND is_hard_gpt IS NULL AND is_custom IS NOT TRUE
		LIMIT 1`
	}
	
	rows, err := d.Connection.Query(query)
	checkErr(err)
	for rows.Next() {
		var name string
		var id int
		var isValid bool

		err = rows.Scan(&id, &name, &isValid)
		skill = models.Skill{
			Id: id,
			Name: name,
			IsValid: isValid,
		}
		fmt.Println("Skill", skill)
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
		SELECT id, translated, is_displayed
		FROM demand
		WHERE translated IS NOT NULL AND type_group IS NULL
		LIMIT 1`
	rows, err := d.Connection.Query(query)
	checkErr(err)
	for rows.Next() {
		var name string
		var id int
		var isValid bool

		err = rows.Scan(&id, &name, &isValid)
		skill = models.Skill{
			Id: id,
			Name: name,
			IsValid: isValid,
		}
	}
	return
}


func (d *Database) UpdateSkillGroup(skill models.Skill) {
	query := fmt.Sprintf(`
		UPDATE demand
		SET type_group='%s'
		WHERE lower(translated) = '%s' `, skill.Group, strings.ToLower(skill.Name))
	d.ExecuteQuery(query)
}