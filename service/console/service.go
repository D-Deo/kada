//Package console 控制台模块.
//	可以在控制台输入指令来控制程序
package console

import (
	"kada/core"
)

const (
	MODULE = "console"
)

var (
	_service *core.Service
)

func init() {
	_service = core.NewService()
	_service.Register(MODULE, NewHandler())
	_service.Start()
}

func Register(cmd string, fun func(...string)) error {
	args := &RegisterArgs{}
	args.Cmd = cmd
	args.Func = fun

	if err := _service.Call(MODULE, "Register", args, nil); err != nil {
		return err
	}
	return nil
}

func Listen() error {
	if err := _service.Call(MODULE, "Listen", nil, nil); err != nil {
		return err
	}
	return nil
}
