package config

type SrvInfoTemplate struct {
	MainVer  string
	ScVer    string
	ModVer   string
	Salt     string
	AllowCmd []string
}

var SrvInfo = SrvInfoTemplate{
	AllowCmd: []string{"list", "say", "tell", "me"},
	Salt:     "Axolotland Gaming Club",
}
