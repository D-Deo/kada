package console

import (
	"bufio"
	"os"
	"strings"

	"github.com/D-Deo/kada.go"
	"github.com/D-Deo/kada.go/log"
)

var cbs map[string]func(...string)

func init() {
	cbs = make(map[string]func(...string))
}

// 注册控制台消息
func Register(cmd string, cb func(...string)) {
	cbs[cmd] = cb
}

// 监听控制台消息
func Listen() {
	log.Signal("[console] wait listening cmd ...")
	reader := bufio.NewReader(os.Stdin)
	for {
		in, _, err := reader.ReadLine()
		if err != nil {
			log.Error("[console] error: %v", err)
			break
		}

		cmds := strings.Split(string(in), " ")
		log.Info("[console] cmd: %v", cmds)
		if cmds[0] == "over" {
			break
		}
		cb, ok := cbs[cmds[0]]
		if !ok {
			log.Warn("[console] no cmd: %s", cmds[0])
			continue
		}
		go func() {
			defer kada.Panic()
			cb(cmds[1:]...)
		}()
	}
}
