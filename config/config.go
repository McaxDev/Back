package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	RconPwd  string
	GptToken string
	JwtKey   []byte
}

var Config = config{
	RconPwd:  "",
	GptToken: "",
	JwtKey:   nil,
}

func ReadConf() error {
	path := "config.yaml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		data, err := yaml.Marshal(&Config)
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
		if err := yaml.Unmarshal(data, &Config); err != nil {
			return err
		}
		return nil
	}
}
