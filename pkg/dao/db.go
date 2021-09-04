package dao

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"mlauth/pkg/conf"
)

func getDb() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", conf.DbDsn)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
