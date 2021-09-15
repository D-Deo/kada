package log

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/issue9/term/colors"
	"github.com/longbozhan/timewriter"
)

type Level int

const (
	_ Level = iota
	LSignal
	LError
	LWarn
	LInfo
	LDebug
	LDump
)

var (
	_option Option
	_logger *log.Logger

	_l = flag.String("l", "debug", "log output level")
	_c = flag.Bool("c", true, "log output console")
)

type Option struct {
	Level   Level
	Console bool
}

func init() {
	flag.Parse()

	switch *_l {
	case "debug":
		_option.Level = LDebug
	case "info":
		_option.Level = LInfo
	case "warn":
		_option.Level = LWarn
	case "error":
		_option.Level = LError
	}
	_option.Console = *_c

	writer := &timewriter.TimeWriter{
		Dir:        "./logs",
		Compress:   true,
		ReserveDay: 30,
	}
	_logger = log.New(writer, "", 0)
}

func output(level Level, output string) {
	ts := time.Now().Format("2006-01-02 15:04:05.999999") //设定时间格式

	switch level {
	case LDebug:
		if _option.Level >= level {
			if _option.Console {
				colors.Println(colors.Cyan, colors.Default, fmt.Sprintf("%s %-26s %s", "[D]", ts, output))
			}
			write("[D]", output)
		}
	case LInfo:
		if _option.Level >= level {
			if _option.Console {
				colors.Println(colors.White, colors.Default, fmt.Sprintf("%s %-26s %s", "[I]", ts, output))
			}
			write("[I]", output)
		}
	case LSignal:
		if _option.Level >= level {
			if _option.Console {
				colors.Println(colors.Green, colors.Default, fmt.Sprintf("%s %-26s %s", "[S]", ts, output))
			}
			write("[S]", output)
		}
	case LWarn:
		if _option.Level >= level {
			if _option.Console {
				colors.Println(colors.Yellow, colors.Default, fmt.Sprintf("%s %-26s %s", "[W]", ts, output))
			}
			write("[W]", output)
		}
	case LError:
		if _option.Level >= level {
			if _option.Console {
				colors.Println(colors.Red, colors.Default, fmt.Sprintf("%s %-26s %s", "[E]", ts, output))
			}
			write("[E]", output)
		}
	case LDump:
		write("[P]", output)
	default:
		break
	}
}

// 写入日志
func write(level string, s string, v ...interface{}) {
	ts := time.Now().Format("15:04:05.999999") //设定时间格式
	var a []interface{}
	a = append(a, ts)
	a = append(a, s)
	a = append(a, v...)
	_logger.SetPrefix(level)
	_logger.Printf(" %-15s %s", a...)
}

// 打印调试级别日志
func Debug(format string, v ...interface{}) {
	output(LDebug, fmt.Sprintf(format, v...))
}

// 打印信息级别日志
func Info(format string, v ...interface{}) {
	output(LInfo, fmt.Sprintf(format, v...))
}

// 打印警告级别日志
func Warn(format string, v ...interface{}) {
	output(LWarn, fmt.Sprintf(format, v...))
}

// 打印错误级别日志
func Error(format string, v ...interface{}) {
	output(LError, fmt.Sprintf(format, v...))
}

// 打印信号级别日志
func Signal(format string, v ...interface{}) {
	output(LSignal, fmt.Sprintf(format, v...))
}

// 打印崩溃级别日志（格式化）
func Dump(format string, v ...interface{}) {
	output(LDump, fmt.Sprintf(format, v...))
}

// 兼容日志打印函数
func Print(v ...interface{}) {
	output(LDebug, fmt.Sprint(v...))
}

// 兼容日志打印函数（格式化）
func Printf(format string, v ...interface{}) {
	Debug(format, v...)
}

// 兼容日志崩溃函数
func Panic(v ...interface{}) {
	log.Panic(v...)
}

// 兼容日志打印函数
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}
