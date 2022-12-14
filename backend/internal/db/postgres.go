package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(host, port, user, password string) (*Postgres, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password),
	)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}
