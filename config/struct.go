package config

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

var dbStructs = []interface{}{
	&Log{},
	&User{},
	&Text{},
	&IPs{},
	&Online{},
	&SysLog{},
}

func AutoMigrate() {
	for _, model := range dbStructs {
		if err := DB.AutoMigrate(model); err != nil {
			fmt.Println("自动迁移失败：" + err.Error())
		}
	}
}

type Log struct {
	ID       uint
	Username string
	Time     time.Time
	Status   int
	Error    string
	Duration time.Duration
	Method   string `gorm:"size:10"`
	Path     string `gorm:"size:30"`
	Source   string `gorm:"size:50"`
	ReqBody  string `gorm:"type:text"`
	ResBody  string `gorm:"type:text"`
	Agent    string `gorm:"type:text"`
}

type User struct {
	gorm.Model
	Username   string `gorm:"unique_index"`
	Admin      bool
	Avatar     string
	Password   string `gorm:"size:50"`
	Gamename   string `gorm:"size:30"`
	Telephone  string `gorm:"size:20"`
	Email      string `gorm:"size:50"`
	Profile    string `gorm:"type:text"`
	WhiteCoin  int
	BlueCoin   int
	GptThreads []GptThread `gorm:"foreignkey:UserID"`
}

type Text struct {
	gorm.Model
	Type    string
	Title   string
	Content string `gorm:"type:text"`
	Author  string
}

type IPs struct {
	ID   uint
	Time time.Time
	Ipv4 string `gorm:"size:20"`
	Ipv6 string `gorm:"size:50"`
}

type Online struct {
	ID   uint
	Time time.Time
	Main int
	Sc   int
	Mod  int
}

type SysLog struct {
	ID      uint
	Time    time.Time
	Level   string `gorm:"size:5"`
	Message string `gorm:"type:text"`
}

type GptThread struct {
	ThreadID   string `gorm:"primaryKey"`
	ThreadName string
	UserID     uint `gorm:"index"`
}
