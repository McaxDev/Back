package util

import (
	"strconv"

	co "github.com/McaxDev/Back/config"
	"github.com/gorcon/rcon"
)

func Rcon(server, command string) (string, error) {
	rconAddr := co.Config.ServerIP + ":" + strconv.Itoa(co.Config.RconPort[server])
	conn, err := rcon.Dial(rconAddr, co.Config.RconPwd)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return conn.Execute(command)
}
