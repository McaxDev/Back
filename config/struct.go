package config

import (
	"time"

	"gorm.io/gorm"
)

var TableList = []interface{}{
	&Log{},
	&User{},
	&AxolotlCoin{},
	&Text{},
	&IPs{},
	&Online{},
	&SystemLog{},
	&GptThread{},
}

type Log struct {
	ID       uint `gorm:"primaryKey;auto_increment"`
	User     User `gorm:"foreignKey:UserID;references:ID"`
	UserID   uint `gorm:"column:user_id;type:uint"`
	Time     time.Time
	Status   int
	Error    string
	Duration time.Duration
	Method   string
	Path     string
	Source   string
	ReqBody  string `gorm:"type:text"`
	ResBody  string `gorm:"type:text"`
	Agent    string `gorm:"type:text"`
}

type User struct {
	ID        uint `gorm:"column:user_id;primaryKey;auto_increment;type:uint"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"column:user_name;unique_index;size:255;not null"`
	Admin     bool   `gorm:"column:admin"`
	Avatar    string `gorm:"column:head"`
	Password  string `gorm:"column:user_pas"`
	Gamename  string `gorm:"column:game_name;size:30;not null"`
	Telephone string `gorm:"column:telephone;size:20"`
	Email     string `gorm:"column:email;size:50"`
}

func (User) TableName() string {
	return "userlist"
}

type AxolotlCoin struct {
	gorm.Model
	User   User `gorm:"foreignKey:UserID;references:ID"`
	UserID uint `gorm:"column:user_id;type:uint"`
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
	AuthorID uint   `gorm:"column:author_id;type:uint"`
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

type SystemLog struct {
	gorm.Model
	Level   string `gorm:"column:level"`
	Message string `gorm:"column:message;type:text"`
}

type GptThread struct {
	gorm.Model
	ThreadID   string
	ThreadName string
	UserID     uint `gorm:"column:user_id;type:uint"`
	User       User `gorm:"foreignKey:UserID;references:ID"`
}
