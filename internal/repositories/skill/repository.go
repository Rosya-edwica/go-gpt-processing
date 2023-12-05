package skill

import (
	"fmt"
	"go-gpt-processing/internal/entities"
	"go-gpt-processing/internal/models"
	"strings"

	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetDuplicates() ([]models.Skill, error) {
	query := `
		SELECT id, demand_name as name, dup_demand_name as dup_name
		FROM demand_duplicate
		WHERE is_duplicate_gpt IS NULL
	`
	rawSkills := make([]entities.Skill, 0)
	err := r.db.Select(&rawSkills, query)
	if err != nil {
		return nil, errors.Wrap(err, "select duplicates")
	}
	return models.NewSkills(rawSkills), nil
}

func (r *Repository) UpdateDuplicate(skill models.Skill) (bool, error) {
	res, err := r.db.Exec(`
		UPDATE demand_duplicate
		SET is_duplicate_gpt = ?
		WHERE id = ?
	`, skill.IsDuplicate, skill.Id)

	if err != nil {
		return false, errors.Wrap(err, "update duplicate")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "rows affected duplicates")
	}
	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r *Repository) GetSkillsWithoutTypeGroup() ([]models.Skill, error) {
	rawSkills := make([]entities.Skill, 0)
	err := r.db.Select(&rawSkills, `
		SELECT id, translated as name
		FROM demand
		WHERE translated IS NOT NULL AND type_group IS NULL
	`)
	if err != nil {
		return nil, errors.Wrap(err, "select group skills")
	}
	return models.NewSkills(rawSkills), nil
}

func (r *Repository) UpdateTypeGroup(skill models.Skill) (bool, error) {
	skill.GroupType = strings.TrimSpace(strings.ToLower(skill.GroupType))
	if skill.GroupType != "навык" && skill.GroupType != "другое" && skill.GroupType != "профессия" {
		return false, errors.New(fmt.Sprintf("Недопустимый тип навыка: '%s'. Допустимые типы: навык, профессия, другое", skill.GroupType))
	}
	res, err := r.db.Exec(`
		UPDATE demand
		SET type_group=?
		WHERE id=?
	`, skill.GroupType, skill.Id)
	if err != nil {
		return false, errors.Wrap(err, "update type group")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "rows affected type group")
	}
	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil

}

func (r *Repository) GetPreparedSkills() ([]models.Skill, error) {
	rawSkills := make([]entities.Skill, 0)
	err := r.db.Select(&rawSkills, `
		SELECT id, name
		FROM demand
		WHERE type_group = 'навык' and is_hard_gpt is true AND is_deleted is false
	`)
	if err != nil {
		return nil, errors.Wrap(err, "select prepared skills")
	}
	return models.NewSkills(rawSkills), nil
}

// Передаешь название колонки и функция вернет все строки, где значение этой колонки равно NULL
func (r *Repository) GetNullableSkillsInColumn(column string) ([]models.Skill, error) {
	rawSkills := make([]entities.Skill, 0)
	err := r.db.Select(&rawSkills, `
		SELECT id, name
		FROM demand
		WHERE %s IS NULL
	`)
	if err != nil {
		return nil, errors.Wrap(err, "select prepared skills")
	}
	return models.NewSkills(rawSkills), nil
}

func (r *Repository) UpdateSkillByColumn(skillID int, value, column string) (bool, error) {
	res, err := r.db.Exec(`
		UPDATE demand_duplicate
		SET ? = ?
		WHERE id = ?
	`, column, value, skillID)

	if err != nil {
		return false, errors.Wrap(err, "update skill by column "+column)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "rows affected by column "+column)
	}
	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

// TODO: Сохранить поднавыки и тесты с помощью  INSERT BATCH
