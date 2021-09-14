package gate

import (
	"kada/core"
	"kada/log"
	"kada/utils/config"
)

const (
	SOCKET_MODE    = "1"
	WEBSOCKET_MODE = "2"
)

var (
	_server core.IServer
)

//Startup 启动服务
func Startup() error {
	log.Info("[Gate] Service Startup ...")

	mode := config.GetWithDef("gate", "mode", SOCKET_MODE)
	switch mode {
	case SOCKET_MODE:
		s := new(Server)
		_server = s
	case WEBSOCKET_MODE:
		s := new(WServer)
		_server = s
	default:
		log.Error("[Gate] UnKnow Mode", mode)
		return core.ErrServer
	}

	if err := _server.Startup(); err != nil {
		return err
	}

	log.Info("[Gate] Service Finish ...")
	return nil
}

//Send 发送数据
func Send(sid string, pid int32, data []byte) error {
	return _server.Send(sid, pid, data)
}

//SendAll 发送数据全体
func SendAll(pid int32, data []byte) error {
	return _server.SendAll(pid, data)
}
