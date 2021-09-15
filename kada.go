package kada

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/D-Deo/kada.go/log"
)

type (
	IService interface {
		Startup() error
		Call(string, string, interface{}, interface{}) error
	}

	IHandler interface {
		Handle(string, int32, []byte) error
	}

	Message struct {
		Handle string
		Action string
		Args   interface{}
		Back   interface{}
	}

	Service struct {
		Recv chan Message
		Send chan error

		Handlers map[string]reflect.Value
	}
)

// 创建服务
func NewService() *Service {
	service := &Service{}
	service.Recv = make(chan Message)
	service.Send = make(chan error)
	service.Handlers = make(map[string]reflect.Value)
	return service
}

// 注册服务
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

	for msg := range o.Recv {
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
			}
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

// 启动程序
func Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	log.Signal("[kada] run ...")
	sig := <-c
	log.Signal("[kada] on signal: %v", sig)

	// 如果一分钟都关不了则强制关闭
	timeout := time.NewTimer(time.Second * time.Duration(60))
	wait := make(chan struct{})
	go func() {
		wait <- struct{}{}
	}()
	select {
	case <-timeout.C:
		log.Panic("[kada] close timeout (signal: %v)", sig)
	case <-wait:
		log.Signal("[kada] close down (signal: %v)", sig)
	}
}

// 捕获异常
func Panic() {
	if err := recover(); err != nil {
		// exeName := os.Args[0] //获取程序名称
		// pid := os.Getpid() //获取进程ID
		now := time.Now() //获取当前时间

		time := now.Format("2006_01-02_15-04-05") //设定时间格式
		filename := fmt.Sprintf("%s.dmp", time)   //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）
		fmt.Println("dump to file", filename)

		if f, e := os.Create(filename); e != nil {
			return
		} else {
			defer f.Close()

			f.WriteString(fmt.Sprintf("%v\r\n", err)) //输出panic信息
			f.WriteString(string(debug.Stack()))      //输出堆栈信息
		}
	}
}
