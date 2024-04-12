package config

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

var dbStructs = []interface{}{
	&Log{},
	&User{},
	&AxolotlCoin{},
	&Text{},
	&IPs{},
	&Online{},
	&SysLog{},
	&GptThread{},
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
	User     User `gorm:"foreignKey:userID;references:ID"`
	UserID   uint `gorm:"column:user_id"`
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
	ID        uint   `gorm:"column:user_id"`
	Username  string `gorm:"column:user_name;unique_index"`
	Admin     bool   `gorm:"column:admin"`
	Avatar    string `gorm:"column:head"`
	Password  string `gorm:"column:user_pas"`
	Gamename  string `gorm:"column:game_name;size:30"`
	Telephone string `gorm:"column:telephone;size:20"`
	Email     string `gorm:"column:email;size:50"`
}

func (User) TableName() string {
	return "userlist"
}

type AxolotlCoin struct {
	gorm.Model
	User   User `gorm:"foreignKey:UserID;references:ID"`
	UserID uint `gorm:"column:user_id"`
	Pearl  int  `gorm:"column:pearl_axolotl_coin"`
	Azure  int  `gorm:"column:azure_axolotl_coin"`
}

func (AxolotlCoin) TableName() string {
	return "axolotl_coin"
}

type Text struct {
	gorm.Model
	Type     string `gorm:"column:type"`
	Title    string `gorm:"column:title"`
	Content  string `gorm:"column:content;type:text"`
	AuthorID uint   `gorm:"column:author_id"`
	User     User   `gorm:"foreignKey:AuthorID;references:ID"`
}

type IPs struct {
	gorm.Model
	Ipv4 string `gorm:"column:ipv4"`
	Ipv6 string `gorm:"column:ipv6"`
}

type Online struct {
	gorm.Model
	Main int `gorm:"column:main"`
	Sc   int `gorm:"column:sc"`
	Mod  int `gorm:"column:mod"`
}

func (Online) TableName() string {
	return "online_info"
}

type SysLog struct {
	gorm.Model
	Level   string `gorm:"column:level"`
	Message string `gorm:"column:message;type:text"`
}

type GptThread struct {
	gorm.Model
	ThreadID   string
	ThreadName string
	UserID     uint `gorm:"index"`
	User       User `gorm:"foreignKey:UserID;references:ID"`
}
