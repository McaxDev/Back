package entity

import (
	"time"

	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	Time     string        `gorm:"type:datetime"`
	Level    string        `gorm:"type:varchar(50)"`
	Status   int           `gorm:"type:int"`
	Error    string        `gorm:"type:varchar(255)"`
	Method   string        `gorm:"type:varchar(10)"`
	Path     string        `gorm:"type:varchar(255)"`
	Source   string        `gorm:"type:varchar(100)"`
	Body     string        `gorm:"type:text"`
	Duration time.Duration `gorm:"type:bigint"`
}

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);unique_index"`
	Password string `gorm:"type:varchar(100)"`
	Admin    int    `gorm:"type:int"`
	Gamename string `gorm:"type:varchar(100)"`
	Avatar   string `gorm:"type:text"`
}

type Text struct {
	gorm.Model
	Type    string `gorm:"type:varchar(255)"`
	Title   string `gorm:"type:varchar(255)"`
	Content string `gorm:"type:text"`
}
