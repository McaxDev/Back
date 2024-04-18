package util

import (
	"encoding/json"
	"fmt"
	"net"

	co "github.com/McaxDev/Back/config"
)

// 代表TCP应答格式的结构体
type tcpResp struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

// 记录TCP连接过程中传输的成功信息
func TcpInfo(conn net.Conn, data any) {
	tcpSend(conn, true, data, nil)
}

// 记录TCP连接过程中传输的失败信息
func TcpError(conn net.Conn, data any, errinfo error) {
	tcpSend(conn, false, data, errinfo)
}

// 给上面两个函数用的函数
func tcpSend(conn net.Conn, success bool, data any, errinfo error) {

	// 将消息结构体序列化为字节切片
	marshed, err := json.Marshal(&tcpResp{Success: success, Data: data})
	if err != nil {
		co.SysLog("ERROR", "将TCP消息序列化失败"+err.Error())
		return
	}

	// 将序列化信息发送给客户端
	_, err = conn.Write(marshed)
	if err != nil {
		co.SysLog("ERROR", "无法发送信息到客户端"+err.Error())
		return
	}

	// 从连接里读取连接者的IP
	ip := conn.RemoteAddr().String()

	// 将消息和错误变为字符串
	datastr, errstr := string(marshed), ErrToStr(errinfo)

	// 将消息输出到数据库
	if err := co.DB.Create(&co.TcpLog{
		IP: ip, Data: datastr,
		Error: errstr, IO: false,
	}).Error; err != nil {
		co.SysLog("ERROR", "无法将信息输出到数据库"+err.Error())
	}

	// 将消息输出到控制台
	if errinfo != nil {
		fmt.Print(fmt.Sprintf("ERROR %s ", errinfo.Error()))
	}
	fmt.Println(fmt.Sprintf("TCP %s %s", ip, datastr))
}

// 记录TCP连接本身可能产生的消息
func ConnInfo(conn net.Conn, level, msg string, errinfo error) {

	// 将错误信息参数变成字符串
	errstr := ErrToStr(errinfo)

	// 从连接里读取连接者的IP
	ip := conn.RemoteAddr().String()

	// 将消息输出到数据库
	if err := co.DB.Create(&co.ConnLog{
		IP: ip, Level: level, Info: msg + errstr,
	}).Error; err != nil {
		co.SysLog("ERROR", "无法将TCP连接日志记录到数据库"+err.Error())
	}

	// 将日志消息输出到控制台
	fmt.Println(fmt.Sprintf("%s %s %s", level, ip, msg+errstr))

	// 如果连接失败，则向客户端返回错误提示
	if level == "ERROR" {
		if _, err := conn.Write([]byte(msg)); err != nil {
			co.SysLog("ERROR", "无法向客户端发送错误信息")
			return
		}
	}
}
