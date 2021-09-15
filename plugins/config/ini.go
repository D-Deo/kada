package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/larspensjo/config"
)

var (
	ErrNoValue = errors.New("config error, no value") //没有对应的值
)

// Ini 配置对象
type Ini struct {
	cfg *config.Config
}

// NewIni 创建新的配置对象
func NewIni() *Ini {
	return &Ini{}
}

func (o *Ini) Load(filename string) error {
	if !strings.Contains(filename, ".ini") {
		filename += ".ini"
	}
	cfg, err := config.ReadDefault(filename)
	if err != nil {
		return fmt.Errorf("[config] load (%s) err: %w", filename, err)
	}
	o.cfg = cfg
	return nil
}

// 获取内容
func (o *Ini) Get(title string, key string) (string, error) {
	if !o.cfg.HasOption(title, key) {
		return "", ErrNoValue
	}

	value, err := o.cfg.String(title, key)
	if err != nil {
		return "", err
	}

	return value, nil
}

// 获取内容，如果不存在则返回默认值
func (o *Ini) GetWithDef(title string, key string, def string) string {
	value, err := o.cfg.String(title, key)
	if err != nil {
		return def
	}
	return value
}
