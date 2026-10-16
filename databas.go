package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUser     = "root"
	dbPass     = "root"
	dbHost     = "127.0.0.1"
	dbPort     = "3306"
	dbName     = "minecraft"
	dbTable    = "accounts"
	colUser    = "username"
	colPass    = "password"
	colSesID   = "session"
	colServer  = "server"
	gameTable  = "data"
	propColumn = "property"
	valColumn  = "value"
)

var db *sql.DB

func initDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB Error:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + dbTable + ` (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		` + colUser + ` VARCHAR(255) NOT NULL,
		` + colPass + ` VARCHAR(255) NOT NULL,
		` + colSesID + ` VARCHAR(255) DEFAULT NULL,
		` + colServer + ` VARCHAR(255) DEFAULT NULL,
		skin_path VARCHAR(255) DEFAULT 'default.png',
		cloak_path VARCHAR(255) DEFAULT NULL
	)`)
	if err != nil {
		log.Fatal("Create accounts error:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + gameTable + ` (
		` + propColumn + ` VARCHAR(255) UNIQUE,
		` + valColumn + ` VARCHAR(255)
	)`)
	if err != nil {
		log.Fatal("Create data error:", err)
	}

	db.Exec("INSERT IGNORE INTO " + gameTable + " (" + propColumn + "," + valColumn + ") VALUES ('latest-game-build','10746'),('launcher-version','13')")

	fmt.Println("DB connected")
}

func getGameInfo(typ string) string {
	var prop string
	switch typ {
	case "build":
		prop = "latest-game-build"
	case "launcher":
		prop = "launcher-version"
	default:
		return ""
	}
	var val string
	err := db.QueryRow("SELECT "+valColumn+" FROM "+gameTable+" WHERE "+propColumn+"=?", prop).Scan(&val)
	if err != nil {
		return ""
	}
	return val
}