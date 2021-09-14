package kada

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

var (
	_service *Service
)

func init() {
	_service = NewService()
}

// 注册服务
func Register(name string, handler interface{}) {
	_service.Register(name, handler)
}

// 启动服务
func Start() {
	_service.Start()
}

// 调用服务
func Call(handle string, action string, args interface{}, back interface{}) error {
	return _service.Call(handle, action, args, back)
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

		f, e := os.Create(filename)
		defer f.Close()
		if e != nil {
			return
		}

		f.WriteString(fmt.Sprintf("%v\r\n", err)) //输出panic信息
		f.WriteString(string(debug.Stack()))      //输出堆栈信息
	}
}
