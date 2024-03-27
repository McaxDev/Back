package config

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var Config = struct {
	RconPwd  string
	GptToken string
	JwtKey   []byte
	AllowCmd []string
	McFont   string
	Sql      string
}{}

var Info = struct {
	MainVer string
	ScVer   string
	ModVer  string
}{}

func Read(config interface{}, path string) error {
	var marshalFunc func(interface{}) ([]byte, error)
	var unmarshalFunc func([]byte, interface{}) error
	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		marshalFunc, unmarshalFunc = yaml.Marshal, yaml.Unmarshal
	case ".json":
		marshalFunc, unmarshalFunc = json.Marshal, json.Unmarshal
	case ".xml":
		marshalFunc, unmarshalFunc = xml.Marshal, xml.Unmarshal
	default:
		return errors.New("此文件扩展类型不支持")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		data, err := marshalFunc(config)
		if err != nil {
			return err
		}
		return os.WriteFile(path, data, 0644)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return unmarshalFunc(data, config)
}
