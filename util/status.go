package util

import (
	co "github.com/McaxDev/Back/config"
	"github.com/mcstatus-io/mcutil/v3"
	"github.com/mcstatus-io/mcutil/v3/response"
)

// 执行基础查询服务器信息的函数
func Status(server string) (*response.BasicQuery, error) {

	// 根据参数获取目标服务器的端口
	port := uint16(co.Config.ServerPort[server])

	// 查询服务器信息
	ctx, cancel := Timeout(5)
	defer cancel()
	resp, err := mcutil.BasicQuery(ctx, co.Config.ServerIP, port)
	if err != nil {
		return nil, err
	}

	// 返回服务器信息
	return resp, nil
}

// 执行完整查询服务器信息的函数
func FullStatus(server string) (*response.FullQuery, error) {

	// 根据参数获取目标服务器端口
	port := uint16(co.Config.ServerPort[server])

	// 查询服务器信息
	ctx, cancel := Timeout(5)
	defer cancel()
	resp, err := mcutil.FullQuery(ctx, co.Config.ServerIP, port)
	if err != nil {
		return nil, err
	}

	// 返回服务器信息
	return resp, nil
}
