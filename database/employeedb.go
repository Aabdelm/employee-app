package employeedb

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type DbMap struct {
	l  *log.Logger    // for logging
	m  map[int]string // for storing key-pair vals of departments
	db *sql.DB        //the database
}

func NewDbMap(l *log.Logger) *DbMap {
	return &DbMap{l: l, m: make(map[int]string)}
}

func (DbMap *DbMap) SetupDb() error {
	cfg := &mysql.Config{
		User:   os.Getenv("DBUSERNAME"),
		Passwd: os.Getenv("DBPASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "employee_management",
	}
	//error handling 1
	var err error
	DbMap.db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	//error handling 2
	if err := DbMap.db.Ping(); err != nil {
		return err
	}

	DbMap.l.Println("Connected Database")
	return nil
}
