package util

import (
	"context"
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/mcstatus-io/mcutil/v3"
	"github.com/mcstatus-io/mcutil/v3/options"
	"github.com/mcstatus-io/mcutil/v3/response"
)

// 执行基础查询服务器信息的函数
func Statusa(server string) (*response.BasicQuery, error) {

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

type QueryFunc[T any] func(ctx context.Context, ip string, port uint16, options ...options.Query) (T, error)

// 查询服务器信息的泛型函数
func Status[T any](server string, query QueryFunc[T]) (T, error) {

	// 获取目标服务器的端口
	port := uint16(co.Config.ServerPort[server])

	// 创建允许超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行查询
	resp, err := query(ctx, co.Config.ServerIP, port)
	if err != nil {
		return *new(T), err
	}
	return resp, nil
}
