# kada.go

[wiki](https://github.com/D-Deo/kada.go/wiki)

## 【Redis】

```go
// 导入模块
import "github.com/D-Deo/kada.go/service/redis"
```
```go
// 注册 redis 连接信息
redis.Set(tag, redis.NewClient(host, port, db, pass))

// 执行 redis 命令（命令参考官方，完全一致）
reply, err := redis.Get(tag).Exec(cmd, key, args)

// 快捷指令（GET）
str = redis.Get(tag).Get(key)

// 快捷指令（HGETALL）
maps = redis.Get(tag).HGetAll(key)
```

## 【MongoDB】
```go
// 导入模块
import "github.com/D-Deo/kada.go/service/mongo"
```
```go
// 连接MongoDB
// params:
//  -- uri: mongo连接字符串, ep: mongo://xxx.xxx.xxx.xxx:xxxx
//  -- db: 数据库
mongo.Connect(uri, db)
```

## 【Thrift】
```go
// 导入模块
import "github.com/D-Deo/kada.go/rpc/thrift"
```
```go
// 创建服务监听对象
handler := &XXXService{}
// 创建服务进程对象
processor := ads.NewXXXServiceProcessor(handler)
// 创建ThriftRPC服务
server := thrift.NewServer()
// 启动服务
server.Start(port, processor, "", true, false, false)
```
