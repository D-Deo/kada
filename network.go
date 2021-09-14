package kada

import (
	"errors"
	"net"

	"golang.org/x/net/websocket"
)

var (
	ErrServer = errors.New("server error, please check error log") //服务端错误，需要查看错误日志
	ErrClient = errors.New("client error, please check warn log")  //客户端逻辑异常
)

type GateRequest struct {
	HOST string `json:"host"`
	PORT string `json:"port"`
}

type GateResponse struct {
	Code uint8  `json:"code"`
	Msg  string `json:"msg"`
}

type Session struct {
	Id     string
	Chan   chan []byte
	Conn   net.Conn
	WSConn *websocket.Conn
}

type IServer interface {
	Startup() error
	Send(string, int32, []byte) error
	SendAll(int32, []byte) error
}

type IService interface {
	Startup() error
	Call(string, string, interface{}, interface{}) error
}

type IHandler interface {
	Handle(string, int32, []byte) error
}
