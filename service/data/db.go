package data

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/D-Deo/kada.go"
	"github.com/D-Deo/kada.go/log"
	"github.com/D-Deo/kada.go/utils/config"

	"database/sql"

	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/mattn/go-adodb"
	_ "github.com/lib/pq"
)

const DB_MAX_WORKER = 100

var _db *DB

var (
	_mutex = sync.Mutex{}
)

//GetDB 获取数据库访问服务
func GetDB() *DB {
	if _db == nil {
		_db = &DB{}
	}
	return _db
}

type DBMessage struct {
	SQL  string
	Args []interface{}
	Dest []interface{}
	ReCh chan error
}

type DBInsert struct {
	SQL  string
	Args []interface{}
}

//DB 数据库访问对象
type DB struct {
	Conn      *sql.DB
	Chan      chan DBMessage
	AsyncChan chan DBInsert

	CurChanNum uint64

	MsgList *list.List
}

//Init 初始化数据库
func (o *DB) Init() error {
	log.Info("[db] service startup ...")

	host := config.GetWithDef(config.Db, config.DbHost, "127.0.0.1")
	port := config.GetWithDef(config.Db, config.DbPort, "3306")
	user := config.GetWithDef(config.Db, config.DbUser, "root")
	pass := config.GetWithDef(config.Db, config.DbPass, "123123")
	name := config.GetWithDef(config.Db, config.DbName, "zydb")

	// source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
	source := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
	log.Info("[db] source: %s", source)

	// db, err := sql.Open("mysql", source)
	db, err := sql.Open("postgres", source)
	// db, err := sql.Open("adodb", "Provider=SQLOLEDB;Data Source=192.168.8.180;Initial Catalog=U3DMDM;User ID=sa;Password=123123")
	if err != nil {
		return err
	}

	//持久化时，不能关闭
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	o.Conn = db
	o.Conn.SetMaxIdleConns(DB_MAX_WORKER)
	o.Conn.SetMaxOpenConns(DB_MAX_WORKER)

	o.CurChanNum = 0
	o.Chan = make(chan DBMessage)
	o.AsyncChan = make(chan DBInsert)

	for i := 0; i < DB_MAX_WORKER; i++ {
		go o.Worker(o.Chan, o.AsyncChan, i)
	}

	log.Signal("[db] service finish and start worker", DB_MAX_WORKER)
	return nil
}

//Execute 执行SQL语句
func (o *DB) Execute(sql string, args ...interface{}) error {
	msg := DBInsert{
		SQL: sql,
	}
	msg.Args = append(msg.Args, args...)
	o.AsyncChan <- msg
	return nil
}

//Worker 工作线程
func (o *DB) Worker(ch <-chan DBMessage, ach <-chan DBInsert, proc int) {
	defer kada.Panic()

	for {
		select {
		case msg := <-ch:
			log.Debug("[db] proc(%d)", proc)
			err := o.Conn.QueryRow(msg.SQL, msg.Args...).Scan(msg.Dest...)
			msg.ReCh <- err
		case msg := <-ach:
			log.Info("[db] async exec proc(%d) sql(%s) args(%v)", proc, msg.SQL, msg.Args)
			stmt, err := o.Conn.Prepare(msg.SQL)
			defer stmt.Close()
			if err != nil {
				log.Panic("[db] worker proc(%d) err(%v) - sql(%s) args(%v)", proc, err, msg.SQL, msg.Args)
				return
			}
			stmt.Exec(msg.Args...)
			stmt.Close()
		}
	}
}
