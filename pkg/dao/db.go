package dao

import (
	"github.com/jmoiron/sqlx"
	"mlauth/pkg/conf"
)

func getDb() (*sqlx.DB, error) {
	db, err := sqlx.Connect(conf.DbDriver, conf.DbSource)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
