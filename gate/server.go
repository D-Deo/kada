package gate

import (
	"fmt"
	"net"
	"sync"

	"github.com/D-Deo/kada.go"
	"github.com/D-Deo/kada.go/log"
	"github.com/D-Deo/kada.go/plugins/config"
	"github.com/D-Deo/kada.go/utils"
)

// var _server *Server

// //GetServer 获取服务
// func GetServer() *Server {
// 	if _server == nil {
// 		_server = &Server{}
// 	}
// 	return _server
// }

//Server TCP服务，监听客户端连接和收发消息
type Server struct {
	Sessions map[string]Session
	Locker   sync.Mutex
}

//Startup 启动服务，监听端口
func (o *Server) Startup() error {
	o.Sessions = make(map[string]Session)
	o.Locker = sync.Mutex{}

	host := config.GetWithDef("gate", "host", "127.0.0.1")
	port := config.GetWithDef("gate", "port", "10000")

	address := fmt.Sprintf("%s:%s", host, port)
	log.Info("[Gate] Address", address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	//启动监听
	go o.Listen(listener)

	log.Info("[Gate] Waiting For Clients ...")
	return nil
}

//Listen 监听客户端连接
func (o *Server) Listen(listener net.Listener) {
	defer kada.Panic()
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		// defer conn.Close()

		if err != nil {
			log.Error("[Gate] Accept", err)
			continue
		}

		sid := conn.RemoteAddr().String()
		session := Session{}
		session.Id = sid
		session.Chan = make(chan []byte)
		session.Conn = conn
		o.Sessions[sid] = session
		log.Info(sid, "Gate Server Connect Success")

		go o.Handle(session)
	}
}

//Handle 处理连接
func (o *Server) Handle(session Session) {
	defer kada.Panic()

	go func(s Session) {
		defer kada.Panic()

		for data := range s.Chan {
			if _, err := s.Conn.Write(data); err != nil {
				log.Error(s.Id, "Gate Server Write", err)
				break
			}
		}

		// for {
		// 	select {
		// 	case data := <-s.Chan:
		// 		if _, err := s.Conn.Write(data); err != nil {
		// 			log.Error(s.Id, "Gate Server Write", err)
		// 			break
		// 		}
		// 	}
		// }
	}(session)

	buffer := make([]byte, 0)
	data := make([]byte, 1024)
	for {
		n, err := session.Conn.Read(data)
		if err != nil {
			log.Warn(session.Id, "Gate Server Connection Error", err)
			// o.Locker.Lock()
			// if _, ok := o.Sessions[session.Id]; ok {
			// 	delete(o.Sessions, session.Id)
			// }
			// o.Locker.Unlock()
			return
		}
		buffer = Depack(session.Id, append(buffer, data[:n]...))
		log.Debug(session.Id, "Gate Server Receive Finish", buffer)
	}
}

//Send 发送数据
func (o *Server) Send(sid string, pid int32, data []byte) error {
	if session, ok := o.Sessions[sid]; ok {
		data = Enpack(pid, data)
		session.Chan <- data
		log.Debug(sid, "Gate Server send pid", pid, "data", utils.PrintBuffer(data))
		return nil
	}
	log.Warn(sid, "Gate Server no found client")
	return nil
}

//SendAll 发送数据全体
func (o *Server) SendAll(pid int32, data []byte) error {
	// for sid, conn := range o.conns {
	// 	if conn != nil {
	// 		data = Enpack(pid, data)
	// 		if _, err := conn.Write(data); err != nil {
	// 			log.Error(sid, "Gate Server write error:", err.Error())
	// 			continue
	// 		}
	// 	}
	// }
	log.Info("Gate Server send all pid", pid, "data", data)
	return nil
}
