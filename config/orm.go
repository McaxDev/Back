package config

import (
	"time"

	"gorm.io/gorm"
)

// 对数据库进行注册
var TableList = []interface{}{
	&Log{},
	&User{},
	&AxolotlCoin{},
	&Text{},
	&IPs{},
	&Online{},
	&SystemLog{},
	&GptThread{},
	&TcpLog{},
	&ConnLog{},
}

// 记录HTTP请求信息
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

// 记录用户列表
type User struct {
	ID        uint `gorm:"column:user_id;primaryKey;auto_increment;type:uint" json:"uid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string      `gorm:"column:user_name;unique_index;size:255;not null" json:"username"`
	Admin     bool        `gorm:"column:admin" json:"admin"`
	Avatar    string      `gorm:"column:head" json:"avatar"`
	Password  string      `gorm:"column:user_pas" json:"status"`
	Gamename  string      `gorm:"column:game_name;size:30;not null" json:"gamename"`
	Telephone string      `gorm:"column:telephone;size:20" json:"telephone"`
	Email     string      `gorm:"column:email;size:50" json:"email"`
	GameAuth  bool        `gorm:"column:game_auth" json:"gameauth"`
	Nonce     string      `gorm:"column:nonce" json:"nonce"`
	Balance   AxolotlCoin `gorm:"foreignKey:UserID;references:ID"`
	Thread    []GptThread `gorm:"foreignKey:UserID;references:ID"`
}

func (User) TableName() string {
	return "userlist"
}

// 记录玩家的币
type AxolotlCoin struct {
	gorm.Model
	UserID uint `gorm:"column:user_id;type:uint;uniqueIndex"`
	Pearl  int  `gorm:"column:pearl_axolotl_coin;default:50"`
	Azure  int  `gorm:"column:azure_axolotl_coin;default:0"`
}

// 为结构体指定在数据库对应的名称
func (AxolotlCoin) TableName() string {
	return "axolotl_coin"
}

// 记录在数据库存储的文本
type Text struct {
	gorm.Model
	Type     string `gorm:"column:type"`
	Title    string `gorm:"column:title"`
	Content  string `gorm:"column:content;type:text"`
	AuthorID uint   `gorm:"column:author_id;type:uint"`
	User     User   `gorm:"foreignKey:AuthorID;references:ID"`
}

// 记录服务器的历史IP地址
type IPs struct {
	gorm.Model
	Ipv4 string `gorm:"column:ipv4"`
	Ipv6 string `gorm:"column:ipv6"`
}

// 记录在线玩家历史
type Online struct {
	gorm.Model
	Main int `gorm:"column:main"`
	Sc   int `gorm:"column:sc"`
	Mod  int `gorm:"column:mod"`
}

func (Online) TableName() string {
	return "online_info"
}

// 记录后端系统日志
type SystemLog struct {
	gorm.Model
	Level   string `gorm:"column:level"`
	Message string `gorm:"column:message;type:text"`
}

// 记录不同的用户对应哪些gpt会话
type GptThread struct {
	gorm.Model
	ThreadID   string
	ThreadName string
	UserID     uint `gorm:"column:user_id;type:uint"`
}

// 记录TCP连接与客户端发送的消息的结构体
type TcpLog struct {
	ID        uint
	CreatedAt time.Time
	IP        string
	Data      string
	Error     string
	IO        bool
}

// 记录TCP连接日志的结构体
type ConnLog struct {
	ID        uint
	CreatedAt time.Time
	IP        string
	Level     string
	Info      string
}
