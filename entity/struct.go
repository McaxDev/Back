package entity

import (
	"gorm.io/gorm"
)

type ErrorLog struct {
	ID      int    `gorm:"primary_key"`
	Time    string `gorm:"type:datatime"`
	Level   string `gorm:"type:varchar(50)"`
	Status  int    `gorm:"type:int"`
	Message string `gorm:"type:text"`
}

type User struct {
	gorm.Model
	UserName string
	UserPas  string
	Admin    int
	GameName string
	Head     string
}
