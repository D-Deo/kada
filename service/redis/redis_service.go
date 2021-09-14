package redis

import (
	"fmt"
	"kada/log"

	"github.com/gomodule/redigo/redis"
)

type RedisService struct {
	Conn        redis.Conn
	PubSubConn  redis.PubSubConn
	PubSubCbMap map[string]func(channel, message string)
}

func NewRedisService() *RedisService {
	return &RedisService{}
}

func (o *RedisService) Startup(ip string, port int32) error {
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}
	o.Conn = conn
	o.PubSubConn = redis.PubSubConn{o.Conn}
	o.PubSubCbMap = make(map[string]func(channel, message string))

	go func() {
		log.Info("Redis", "服务启动 ...")
		for {
			switch res := o.PubSubConn.Receive().(type) {
			case redis.Message:
				log.Debug("Redis", "收到消息", res.Channel, (string)(res.Data))
				o.PubSubCbMap[res.Channel](res.Channel, (string)(res.Data))
			case redis.Subscription:
				log.Debug("Redis", "订阅消息", res.Channel, res.Kind, res.Count)
			case error:
				log.Error("Redis", "订阅消息失败...")
				continue
			}
		}
	}()

	return nil
}

func (o *RedisService) Subscribe(channel interface{}, cb func(channel, message string)) error {
	if err := o.PubSubConn.Subscribe(channel); err != nil {
		return err
	}
	o.PubSubCbMap[channel.(string)] = cb
	return nil
}
