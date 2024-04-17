package config

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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
	if err := Read(&SrvInfo, "srvinfo.json"); err != nil {
		log.Fatalf("重新加载信息文件失败：%v", err)
	}
}

// 将系统日志输出到控制台和数据库
func SysLog(level string, mes string) {

	// 将日志输出到数据库
	syslog := SystemLog{
		Level:   level,
		Message: mes,
	}
	if dbErr := DB.Create(&syslog).Error; dbErr != nil {
		log.Println("将日志存储到数据库失败：" + dbErr.Error())
	}

	// 将日志输出到控制台
	logprinted := "系统日志：" + level + mes
	if level == "FATAL" {
		log.Fatal(logprinted)
	}
	fmt.Println(logprinted)
}
