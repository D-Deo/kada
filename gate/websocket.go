package gate

import (
	"net/http"

	"github.com/D-Deo/kada.go/log"

	"golang.org/x/net/websocket"
)

type WebSocektMessage func(string, string)
type WebSocketBinary func(string, []byte)

//WebSocket
type WebSocket struct {
	Sessions       map[string]Session
	BinaryHandler  WebSocketBinary
	MessageHandler WebSocektMessage
}

func NewWebSocket() *WebSocket {
	ws := &WebSocket{}
	ws.Sessions = make(map[string]Session)
	return ws
}

func (o *WebSocket) Connect(port string, path string) {
	addr := ":" + port
	http.Handle(path, websocket.Handler(o.Handler))
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error("WebSocket", err)
	}
}

func (o *WebSocket) Handler(ws *websocket.Conn) {
	defer ws.Close()

	sid := ws.RemoteAddr().String()
	session := Session{}
	session.Id = sid
	session.WSConn = ws
	o.Sessions[sid] = session
	log.Info("WebSocket", sid, "Connect Success")

	for {
		var data []byte

		if err := websocket.Message.Receive(ws, &data); err != nil {
			log.Warn("WebSocket", session.Id, "Receive Error", err)
			return
		}

		o.BinaryHandler(session.Id, data)
	}
}

func (o *WebSocket) OnMessage(handler WebSocektMessage) {
	o.MessageHandler = handler
}

func (o *WebSocket) OnBinary(handler WebSocketBinary) {
	o.BinaryHandler = handler
}

func (o *WebSocket) Send(sid string, data []byte) error {
	return websocket.Message.Send(o.Sessions[sid].WSConn, data)
}
