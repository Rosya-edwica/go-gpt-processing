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

func (d *Database) SaveSubskills(skill models.Skill) {
	query := fmt.Sprintf(`
		INSERT IGNORE INTO subskills(skill_id, subskills)
		VALUES(%d, '%s')`, skill.Id, strings.Join(skill.SubSkills, "|"))
	d.ExecuteQuery(query)
}
