package mongo

import (
	"log"
)

var (
	clients map[string]*Client
)

func init() {
	clients = make(map[string]*Client)
}

// 注册
func Set(name string, client *Client) {
	clients[name] = client
}

// 获取
func Get(name string) *Client {
	client, ok := clients[name]
	if !ok {
		log.Panicf("[redis] no found redis client: %s", name)
	}
	return client
}
