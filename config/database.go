package config

import (
	"database/sql"
	"errors"
	"time"
)

var DB *sql.DB

func ReadDB() error {
	sqlinfo, ok := Conf["sql"].(string)
	if !ok {
		return errors.New("Assertion failed")
	}
	if dbtmp, err := sql.Open("mysql", sqlinfo); err != nil {
		return err
	} else {
		DB = dbtmp
	}
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)
	if err := DB.Ping(); err != nil {
		return err
	}
	return nil
}
