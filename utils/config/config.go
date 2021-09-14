package config

import (
	"strconv"
)

const (
	Db     = "db"
	DbHost = "host"
	DbPort = "host"
	DbUser = "user"
	DbPass = "pass"
	DbName = "name"
)

var (
	ini *Ini
)

// Load 加载配置文件
func Load(filename string) error {
	ini = NewIni()
	if err := ini.Load(filename); err != nil {
		return err
	}
	return nil
}

// 获取内容
func Get(title string, key string) (string, error) {
	return ini.Get(title, key)
}

// 获取内容，如果不存在则返回默认值
func GetWithDef(title string, key string, def string) string {
	return ini.GetWithDef(title, key, def)
}

// 转 int 输出
func ToInt(str string, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str)
}
