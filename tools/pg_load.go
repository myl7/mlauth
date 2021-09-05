package main

import (
	"flag"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"mlauth/pkg/conf"
	"os"
)

func main() {
	f := flag.String("sql", "", "SQL file to be loaded")
	flag.Parse()
	if *f == "" {
		log.Fatalln("No SQL file is provided")
	}

	db, err := sqlx.Connect("postgres", conf.DbDsn)
	if err != nil {
		log.Fatalln(err.Error())
	}

	c, err := os.ReadFile(*f)
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = db.Exec(string(c))
	if err != nil {
		log.Fatalln(err.Error())
	}
}
