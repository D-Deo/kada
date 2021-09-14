package core

import (
	"log"
	"reflect"
)

type Message struct {
	Handle string
	Action string
	Args   interface{}
	Back   interface{}
}

type Service struct {
	Recv chan Message
	Send chan error
	
	Handlers map[string]reflect.Value
}

func NewService() *Service {
	service := &Service{}
	service.Recv = make(chan Message)
	service.Send = make(chan error)
	service.Handlers = make(map[string]reflect.Value)
	return service
}

func (o *Service) Register(name string, handler interface{}) {
	o.Handlers[name] = reflect.ValueOf(handler)
}

// 启动服务
func (o *Service) Start() {
	go o.Handle()
}

// 控制服务
func (o *Service) Handle() {
	defer Panic()
	
	for {
		select {
		case msg, ok := <-o.Recv:
			if ok {
				if handler, ok := o.Handlers[msg.Handle]; ok {
					if action := handler.MethodByName(msg.Action); action.IsValid() {
						args := reflect.ValueOf(msg.Args)
						back := reflect.ValueOf(msg.Back)
						rest := action.Call([]reflect.Value{args, back})
						var err error
						if !rest[0].IsNil() {
							err = rest[0].Interface().(error)
						}
						o.Send <- err
						break
					}
				}
			}
			log.Panic("[service] recv error")
			continue
		}
	}
}

// 调用服务
func (o *Service) Call(handle string, action string, args interface{}, back interface{}) error {
	if args == nil {
		args = new(int)
	}
	
	if back == nil {
		back = new(int)
	}
	
	msg := Message{handle, action, args, back}
	o.Recv <- msg
	
	err := <-o.Send
	return err
}
