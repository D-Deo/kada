package console

import (
	"bufio"
	"os"
	"strings"

	"github.com/D-Deo/kada.go"
	"github.com/D-Deo/kada.go/log"
)

type RegisterArgs struct {
	Cmd  string
	Func func(...string)
}

func NewHandler() *Handler {
	handler := &Handler{}
	handler.FuncMap = make(map[string]func(...string))
	return handler
}

type Handler struct {
	FuncMap map[string]func(...string)
}

// 启动控制台
func (o *Handler) Start(args *int, back *int) error {
	go o.Listen(args, back)
	return nil
}

// 注册控制台消息
func (o *Handler) Register(args *RegisterArgs, back *int) error {
	o.FuncMap[args.Cmd] = args.Func
	return nil
}

// 监听控制台消息
func (o *Handler) Listen(args *int, back *int) error {
	defer kada.Panic()

	log.Signal("[console] wait listening cmd ...")

	reader := bufio.NewReader(os.Stdin)

	for {
		cmd, _, err := reader.ReadLine()
		if err != nil {
			return err
		}

		cmds := strings.Split(string(cmd), " ")
		log.Info("[console] cmd: %v", cmds)

		if cmds[0] == "over" {
			return nil
		}

		fun, ok := o.FuncMap[cmds[0]]
		if !ok {
			log.Warn("[console] no cmd: %s", cmds[0])
			continue
		}
		go fun(cmds[1:]...)
	}
}
