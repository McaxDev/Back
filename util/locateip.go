package util

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

// 将IP地址字符串转换为城市信息，接受c.ClientIP()
func Locateip(clientip string) (string, error) {

	// 打开IP地址数据库
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		return "", err
	}
	defer db.Close()

	// 根据IP地址查询城市
	record, err := db.City(net.ParseIP(clientip))
	if err != nil {
		return "", err
	}

	// 将城市输出为中文格式并返回
	var location string
	if len(record.Subdivisions) > 0 {
		location = record.Subdivisions[0].Names["zh-CN"] + record.City.Names["zh-CN"]
	}
	if location == "" {
		location = "未知"
	}
	return location, nil
}
