package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
)

// 代表数据包信息的结构体
type TcpJson struct {
	Path string `json:"path"`
	Data any    `json:"data"`
}

// 监听tcp连接的路由
func TcpListener() {

	// 创建tcp服务器对象
	ln, err := net.Listen("tcp", ":"+co.Config.TcpPort)
	if err != nil {
		co.SysLog("ERROR", "创建tcp连接的对象失败："+err.Error())
		return
	}
	defer ln.Close()

	// 使用无限循环持续监听连接
	for {
		// 接受新的连接
		conn, err := ln.Accept()
		if err != nil {
			co.SysLog("INFO", "接受连接失败："+err.Error())
			continue
		}

		// 将连接传递给路由器处理
		go TcpRouter(conn)
	}
}

// 分发新连接的路由器
func TcpRouter(conn net.Conn) {
	defer conn.Close()

	// 使用循环持续读取数据直到连接关闭或发生错误
	for {
		// 将tcp数据读取到缓冲区里
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, conn); err != nil {
			if err != io.EOF {
				util.ConnInfo(conn, "ERROR", "无法读取tcp连接", err)
			}
			return
		}

		data := buf.Bytes()

		if len(data) == 0 {
			continue
		}

		// 将数据包json数据反序列化到结构体
		var dataJson TcpJson
		if err := json.Unmarshal(data, &dataJson); err != nil {
			util.ConnInfo(conn, "ERROR", "无法反序列化数据包", err)
			continue
		}

		switch dataJson.Path {
		case "chat":
			TcpChat(dataJson.Data)
		default:
			util.TcpInfo(conn, "没有这个路由")
		}
	}
}

func TcpChat(req any) {
}
