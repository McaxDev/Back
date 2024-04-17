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
	Ports      Server
	RconPorts  Server
	ServerPath Server
	SMTPConfig struct {
		Srv  string
		Port string
		Mail string
		Pwd  string
	}
	AssistantID struct {
		Gpt3 string
		Axo  string
	}
}

type Server struct {
	Main string
	Sc   string
	Mod  string
}

var SrvConf = make(map[string]map[string]string)

var Config = ConfigTemplate{
	McFont:    "/usr/share/fonts/opentype/axo/mc.ttf",
	Sql:       "backend:backend@tcp(localhost:3306)/backend?charset=utf8mb4&parseTime=True&loc=Local",
	ServerIP:  "192.168.50.38",
	BackPort:  "1314",
	ProxyAddr: "http://127.0.0.1:7890",
	Ports: Server{
		Main: "8000",
		Sc:   "8001",
		Mod:  "8002",
	},
	ServerPath: Server{
		Main: "/srv/main/",
		Sc:   "/srv/sc/",
		Mod:  "/srv/mod/",
	},
}
