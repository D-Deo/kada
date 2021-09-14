package gate

import (
	"github.com/D-Deo/kada.go"
	"github.com/D-Deo/kada.go/log"
)

var (
	_handlers = make(map[int32]*ControlHandler)
)

// var _control *Control

// func GetControl() *Control {
// 	if _control == nil {
// 		_control = &Control{}
// 		_control.handlers = make(map[int32]*ControlHandler)
// 	}
// 	return _control
// }

// type Control struct {
// 	handlers map[int32]*ControlHandler
// }

type ControlHandler struct {
	Handle  string
	Action  string
	Service *kada.IService
}

type GateMessage struct {
	Sid  string
	Head int32
	Data []byte
}

//Call 请求控制器
func Call(sid string, head int32, data []byte) {
	handler, ok := _handlers[head]
	if !ok {
		log.Warn("no handler head:", head)
		return
	}

	args := &GateMessage{
		Sid:  sid,
		Head: head,
		Data: data,
	}
	if err := (*handler.Service).Call(handler.Handle, handler.Action, args, nil); err != nil {
		log.Error(sid, head, "server handle error:", err)
		return
	}

	log.Info(sid, head, "server handle success")
}

//Bind 绑定控制器
func Bind(id int32, handle string, action string, service kada.IService) {
	handler := &ControlHandler{
		Handle:  handle,
		Action:  action,
		Service: &service,
	}
	_handlers[id] = handler
}
