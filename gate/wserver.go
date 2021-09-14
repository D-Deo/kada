package gate

import (
	"kada/core"
	"kada/log"
	"kada/utils/config"
	"net/http"

	"golang.org/x/net/websocket"
)

//WServer WebSocket服务端实现
type WServer struct {
	Sessions map[string]core.Session
}

//Startup 启动服务，监听端口
func (o *WServer) Startup() error {
	o.Sessions = make(map[string]core.Session)

	port := config.GetWithDef("gate", "port", "10000")
	log.Info("[Gate] WS Listen Port", port)

	http.Handle("/", websocket.Handler(o.Handle))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return err
	}

	log.Info("[Gate] Waiting For Clients ...")
	return nil
}

func (o *WServer) Handle(ws *websocket.Conn) {
	defer ws.Close()

	sid := ws.RemoteAddr().String()
	session := core.Session{}
	session.Id = sid
	// session.Chan = make(chan []byte)
	session.WSConn = ws
	o.Sessions[sid] = session
	log.Info(sid, "[Gate] Connect Success")

	buffer := make([]byte, 0)
	for {
		var data []byte

		if err := websocket.Message.Receive(ws, &data); err != nil {
			log.Warn(session.Id, "[Gate] Receive Error", err)
			return
		}

		buffer = Depack(session.Id, append(buffer, data...))
		log.Debug(session.Id, "[Gate] Receive Finish", buffer)
	}
}

//Send 发送数据
func (o *WServer) Send(sid string, pid int32, data []byte) error {
	if session, ok := o.Sessions[sid]; ok {
		data = Enpack(pid, data)
		if err := websocket.Message.Send(session.WSConn, data); err != nil {
			log.Panic("[Gate] send", err)
			return err
		}
		log.Debug(sid, "[Gate] send pid", pid, "data", core.PrintBuffer(data))
		return nil
	}
	log.Warn(sid, "[Gate] no found client")
	return nil
}

//SendAll 发送数据全体
func (o *WServer) SendAll(pid int32, data []byte) error {
	return nil
}
