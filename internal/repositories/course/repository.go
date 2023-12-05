package course

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

func (r *Repository) getCoursesByQuery(query string) ([]models.Course, error) {
	rawCourses := make([]entities.Course, 0)
	err := r.db.Select(&rawCourses, query)
	if err != nil {
		return nil, errors.Wrap(err, "course-getting")
	}
	return models.NewCourses(rawCourses), nil
}

func (r *Repository) GetAllCoursesWithoutGPTSkills() ([]models.Course, error) {
	query := `
		SELECT id, name
		FROM demand
		WHERE is_public = 1
		AND id NOT IN (
			SELECT DISTINCT(course_id) FROM demand_to_course WHERE is_chatgpt IS TRUE
		)
	`
	return r.getCoursesByQuery(query)
}

func (r *Repository) GetOpeneduCoursesWithoutGPTSkills() ([]models.Course, error) {
	query := `
		SELECT course.id as id, course.name as name  FROM course
		LEFT JOIN company ON company.id = company_id
		WHERE company.id = 2077 AND course.id NOT IN (
			SELECT course_id FROM demand_to_course WHERE is_chatgpt IS TRUE
		)
	`
	return r.getCoursesByQuery(query)
}

// TODO: Протестить в первую очередь
func (r *Repository) GetSkillIdByName(name string) (int64, error) {
	var skillId int64
	query, args, err := sqlx.In("SELECT id FROM demand WHERE LOWER(name) = ? ORDER BY is_deleted ASC LIMIT 1", name)
	if err != nil {
		return 0, errors.Wrap(err, "build query in")
	}
	err = r.db.Select(&skillId, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "select skill by name")
	}
	return skillId, err
}

func (r *Repository) ConnectCourseWithSkills(courseID int, skillsIDs []int64) error {
	if len(skillsIDs) == 0 {
		return nil
	}
	valuesQuery := make([]string, 0, len(skillsIDs))
	valuesArgs := make([]interface{}, 0, len(skillsIDs))
	for _, skillID := range skillsIDs {
		valuesQuery = append(valuesQuery, fmt.Sprintf("(?, ?, ?)"))
		valuesArgs = append(valuesArgs, courseID, skillID, true)
	}
	query := fmt.Sprintf(`
		INSERT INTO demand_to_course(course_id, demand_id, is_chatgpt) VALUES %s
	`, strings.Join(valuesQuery, ", "))

	_, err := r.db.Exec(query, valuesArgs...)
	if err != nil {
		return errors.Wrap(err, "connecting course with skills")
	}
	return nil
}
