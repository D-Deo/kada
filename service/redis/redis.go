package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Client struct {
	Pool *redis.Pool
}

func NewClient(host string, port int, db int, password string) *Client {
	client := new(Client)
	client.Pool = &redis.Pool{
		MaxIdle:     256,
		MaxActive:   0,
		IdleTimeout: time.Duration(120),
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", host, port),
				redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialDatabase(db),
				redis.DialPassword(password),
			)
		},
	}
	return client
}

func (o *Client) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	conn := o.Pool.Get()
	if err := conn.Err(); err != nil {
		return nil, fmt.Errorf("[redis] conn get err: %w", err)
	}
	defer conn.Close()
	params := make([]interface{}, 0)
	params = append(params, key)

	if len(args) > 0 {
		for _, v := range args {
			params = append(params, v)
		}
	}
	reply, err := conn.Do(cmd, params...)
	if err != nil {
		return nil, fmt.Errorf("[redis] conn do err: %w", err)
	}
	return reply, nil
}

func (o *Client) Get(key string) string {
	ret, err := redis.String(o.Exec("GET", key))
	if err != nil {
		log.Printf("[redis] get %v err: %v", key, err)
		return ""
	}
	return ret
}

func (o *Client) HGetAll(key string) map[string]string {
	ret, err := redis.StringMap(o.Exec("HGETALL", key))
	if err != nil {
		log.Printf("[redis] hgetall %s err: %v", key, err)
		return nil
	}
	return ret
}
