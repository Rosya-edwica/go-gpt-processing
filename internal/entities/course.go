package entities

type Course struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}
