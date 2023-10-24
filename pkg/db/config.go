package db

import (
	"database/sql"
	"fmt"
	"go-gpt-processing/pkg/logger"
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

const logPrefix = "database: "

func (d *Database) Connect() {
	connectionUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.User, d.Password, d.Host, d.Port, d.Name)
	connection, err := sql.Open("mysql", connectionUrl)
	checkErr(err)
	d.Connection = connection
	logger.LogWarning.Println(logPrefix + "Открыли соединение")
}

func (d *Database) Close() {
	d.Connection.Close()
	logger.LogWarning.Println(logPrefix + "Закрыли соединение")
}

func checkErr(err error) {
	if err != nil {
		logger.LogError.Fatal(logPrefix + err.Error())
	}
}

func (d *Database) ExecuteQuery(query string) {
	tx, _ := d.Connection.Begin()
	_, err := d.Connection.Exec(query)
	checkErr(err)
	tx.Commit()
}

func convertArrayToSQLString(items []string) (result string) {
	result = strings.Join(items, "|")
	result = strings.ReplaceAll(result, ".", "")
	result = strings.ToLower(result)
	result = strings.ReplaceAll(result, "'", "`")
	return
}
