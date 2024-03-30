package config

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var Config = struct {
	RconPwd  string
	GptToken string
	JwtKey   string
	AllowCmd []string
	McFont   string
	Sql      string
	Salt     string
	ServerIP string
	Port
	BackPort string
}{
	AllowCmd: []string{"list", "say", "tell", "me"},
	McFont:   "/usr/share/fonts/opentype/axo/mc.ttf",
	Sql:      "backend:backend@tcp(localhost:3306)/backend?charset=utf8mb4&parseTime=True&loc=Local",
	Salt:     "Axolotland Gaming Club",
	ServerIP: "192.168.50.38",
	BackPort: "1314",
}

var SrvInfo = struct {
	MainVer string
	ScVer   string
	ModVer  string
}{}

func jsonMarshal(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}

func xmlMarshal(v interface{}) ([]byte, error) {
	return xml.MarshalIndent(v, "", "    ")
}

func Read(config interface{}, path string) error {
	var marshalFunc func(interface{}) ([]byte, error)
	var unmarshalFunc func([]byte, interface{}) error
	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		marshalFunc, unmarshalFunc = yaml.Marshal, yaml.Unmarshal
	case ".json":
		marshalFunc, unmarshalFunc = jsonMarshal, json.Unmarshal
	case ".xml":
		marshalFunc, unmarshalFunc = xmlMarshal, xml.Unmarshal
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

func LoadConfig() {
	if err := Read(&Config, "config.yaml"); err != nil {
		log.Fatalf("重新加载配置文件失败：%v", err)
	}
	if err := Read(&SrvInfo, "info.json"); err != nil {
		log.Fatalf("重新加载信息文件失败：%v", err)
	}
}

func Init() {
	LoadConfig()
	if err := ReadDB(); err != nil {
		log.Fatal("读取数据库失败：", err)
	}
}
