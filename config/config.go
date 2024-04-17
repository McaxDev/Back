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
	ServerPort map[string]int
	RconPort   map[string]int
	ServerPath map[string]string
	AsstID     map[string]string
	SMTPConfig map[string]string
}

var Config = ConfigTemplate{
	McFont:    "/usr/share/fonts/opentype/axo/mc.ttf",
	Sql:       "backend:backend@tcp(localhost:3306)/backend?charset=utf8mb4&parseTime=True&loc=Local",
	ServerIP:  "192.168.50.38",
	BackPort:  "1314",
	ProxyAddr: "http://127.0.0.1:7890",
	ServerPort: map[string]int{
		"main": 25565,
		"sc":   25566,
		"mod":  25564,
	},
	RconPort: map[string]int{
		"main": 25575,
		"sc":   25576,
		"mod":  25574,
	},
	ServerPath: map[string]string{
		"main": "/srv/main/",
		"sc":   "/srv/sc/",
		"mod":  "/srv/mod/",
	},
	AsstID: map[string]string{
		"GPT3.5": "",
		"GPT4":   "",
		"HELPER": "",
	},
	SMTPConfig: map[string]string{
		"server":   "smtp.163.com",
		"port":     "25",
		"mail":     "axolotland@163.com",
		"password": "",
	},
}
