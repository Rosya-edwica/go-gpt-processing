package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Host       string
	Port       string
	User       string
	Password   string
	Name       string
	Connection *sql.DB
}

func (d *Database) Connect() {
	connectionUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.User, d.Password, d.Host, d.Port, d.Name)
	connection, err := sql.Open("mysql", connectionUrl)
	checkErr(err)
	d.Connection = connection
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (d *Database) ExecuteQuery(query string) {
	tx, _ := d.Connection.Begin()
	_, err := d.Connection.Exec(query)
	checkErr(err)
	tx.Commit()
}

func (d *Database) Close() {
	d.Connection.Close()
}

func convertArrayToSQLString(items []string) (result string) {
	result = strings.Join(items, "|")
	result = strings.ReplaceAll(result, ".", "")
	result = strings.ToLower(result)
	result = strings.ReplaceAll(result, "'", "`")
	return
}
