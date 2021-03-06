package db

import (
	"database/sql"
	"git.oschina.net/pbaapp/goyaf"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

//表结构
type Adapter struct {
	config map[string]string //连接配置
	db     *sql.DB
}

//实例化驱动
func NewAdapter(config map[string]string) *Adapter {
	key := config["driver"] + "-" + config["host"] + "-" + config["port"] +
		"-" + config["database"] + "-" + config["username"] +
		"-" + config["password"] + "-" + config["charset"]
	goyaf.Debug(key)
	goyaf.Debug(adapters)
	adapter, ok := adapters[key]
	if ok {
		return adapter
	}

	if config["driver"] == "mysql" {
		adapter := &Adapter{}

		host := config["host"]
		port := config["port"]
		database := config["database"]
		username := config["username"]
		password := config["password"]
		charset := config["charset"]

		connectString := username + ":" + password +
			"@tcp(" + host + ":" + port + ")/" + database + "?charset=" + charset

		goyaf.Debug("connect db string:", connectString)
		var err error
		adapter.db, err = sql.Open("mysql", connectString)
		if err != nil {
			panic(err)
		}
		adapter.db.Ping()
		maxConn, ok := config["maxconn"]
		if ok {
			mc, _ := strconv.Atoi(maxConn)
			goyaf.Debug("max open conn", mc)
			adapter.db.SetMaxOpenConns(mc)
		}
		maxIdleConn, ok := config["maxidleconn"]
		if ok {
			mic, _ := strconv.Atoi(maxIdleConn)
			goyaf.Debug("max idle conn", mic)
			adapter.db.SetMaxIdleConns(mic)
		}

		adapters[key] = adapter
	}

	return adapters[key]
}

func (this *Adapter) Begin() (*Transaction, error) {
	tx := &Transaction{}

	t, err := this.db.Begin()
	if err != nil {
		return tx, err
	}
	tx.Tx = t

	return tx, nil
}

func (this *Adapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := this.db.Query(query, args...)
	if err != nil {
		goyaf.ErrorLog.Println("query error:", err, ", try open conn")
		for i := 0; i < 10; i++ {
			_, e := this.db.Query("select 1")
			if e == nil {
				break
			}
		}
		rows, err = this.db.Query(query, args...)
	}
	return rows, err
}

func (this *Adapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := this.db.Exec(query, args...)
	return result, err
}

func (this *Adapter) Prepare(query string) (*sql.Stmt, error) {
	stmt, err := this.db.Prepare(query)
	return stmt, err
}
