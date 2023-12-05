package entities

type Skill struct {
	Id            int    `db:"id"`
	Name          string `db:"name"`
	DuplicateName string `db:"dup_name"`
	IsDuplicate   bool   `db:"is_duplicate"`
	IsValid       bool   `db:"is_deleted"`
	GroupType     string `db:"type_group"`
	// SubSkills     []string `db:""`
	// Description   string
}
