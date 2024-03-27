package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Config = struct {
	RconPwd  string
	GptToken string
	JwtKey   []byte
	AllowCmd []string
	McFont   string
}{}

var Info = struct {
	MainVer string
	ScVer   string
	ModVer  string
}{}

func Read(config interface{}, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		data, err := yaml.Marshal(config)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return err
		}
		return nil
	} else {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(data, config); err != nil {
			return err
		}
		return nil
	}
}
