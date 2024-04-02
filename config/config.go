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
	Ports      Server
	ServerPath Server
}

type Server struct {
	Main string `json:"main"`
	Sc   string `json:"sc"`
	Mod  string `json:"mod"`
}

var SrvConf = make(map[string]map[string]string)

var Config = ConfigTemplate{
	McFont:   "/usr/share/fonts/opentype/axo/mc.ttf",
	Sql:      "backend:backend@tcp(localhost:3306)/backend?charset=utf8mb4&parseTime=True&loc=Local",
	ServerIP: "192.168.50.38",
	BackPort: "1314",
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
