package driver

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func ConnectToSQL() *sql.DB {
	cfg := mysql.Config{
		User:   "sahil",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "CarDealership",
	}

	// get a database handle
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected!")

	return db
}
