package entity

import (
	"gorm.io/gorm"
)

type Log struct {
	ID      int    `gorm:"primary_key"`
	Time    string `gorm:"type:datatime"`
	Level   string `gorm:"type:varchar(50)"`
	Status  int    `gorm:"type:int"`
	Message string `gorm:"type:text"`
	Error   string `gorm:"type:varchar(100)"`
}

type User struct {
	gorm.Model
	UserName string
	UserPas  string
	Admin    int
	GameName string
	Head     string
}
