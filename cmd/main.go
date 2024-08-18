package main

import (
	"database/sql"
	"fmt"
	"github.com/JasonAcar/ecommerce/cmd/api"
	"github.com/JasonAcar/ecommerce/config"
	"github.com/JasonAcar/ecommerce/db"
	"github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	d, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	initStorage(d)

	s := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), d)
	log.Fatal(s.Run())
}

func initStorage(d *sql.DB) {
	err := d.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected!")
}
