package config

type ConfigTemplate struct {
	RconPwd    string
	GptToken   string
	JwtKey     string
	McFont     string
	Sql        string
	ServerIP   string
	BackPort   string
	AvatarPath string
	ProxyAddr  string
	Server     map[string]map[string]string
	SMTPConfig struct {
		Srv  string
		Port string
		Mail string
		Pwd  string
	}
	AsstID map[string]string
}

var SrvConf = make(map[string]map[string]string)

var Config = ConfigTemplate{
	McFont:    "/usr/share/fonts/opentype/axo/mc.ttf",
	Sql:       "backend:backend@tcp(localhost:3306)/backend?charset=utf8mb4&parseTime=True&loc=Local",
	ServerIP:  "192.168.50.38",
	BackPort:  "1314",
	ProxyAddr: "http://127.0.0.1:7890",
	Server: map[string]map[string]string{
		"port":     serversconf,
		"rconport": serversconf,
		"path":     serversconf,
	},
	AsstID: map[string]string{
		"GPT3.5": "",
		"GPT4":   "",
		"HELPER": "",
	},
}

var serversconf = map[string]string{
	"main": "",
	"sc":   "",
	"mod":  "",
}
