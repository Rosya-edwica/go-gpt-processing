package mysql

import (
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
)

func New(cfg Config) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB)
	conn, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		return nil, errors.Wrap(err, "mysql-connection")
	}
	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "mysql-ping-failed")
	}

	return conn, nil
}

func Close(conn *sqlx.DB) {
	conn.Close()
}
