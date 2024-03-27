package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	UserName string
	UserPas  string
	Admin    int
	GameName string
	Head     string
}

func ReadDB() error {
	var err error
	DB, err = gorm.Open(mysql.Open(Config.Sql), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
