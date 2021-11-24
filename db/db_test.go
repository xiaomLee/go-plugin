package db

import (
	"context"
	"testing"
	"time"
)

func TestAddDB_MYSQL(t *testing.T) {
	addDB(MYSQL)

	db := GetDB("user")
	if err := db.Exec("select 3*3").Error; err != nil {
		t.Error(err)
	}
	MustGetDB("not-exist-db")
}

func TestAddDB_PLSQL(t *testing.T) {
	addDB(POSTGRES)

	db := GetDB("test")
	if err := db.Exec("select 3*3").Error; err != nil {
		t.Error(err)
	}
	MustGetDB("not-exist-db")
}

func addDB(dbType string, opts ...Option) {
	switch dbType {
	case MYSQL:
		AddDB("mysql", "user", "root:root@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local", opts...)
	case POSTGRES:
		AddDB("postgres", "test", "host=10.152.239.200 user=db password=db dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai", opts...)

	}
}

func TestDBOptions(t *testing.T) {
	addDB(POSTGRES,
		MaxConn(1),
		MaxOpenConn(1),
		IdleConn(1),
		MaxLeftTime(1),
	)

	instance := MustGetDB("test")
	db, err := instance.DB.DB()
	if err != nil {
		t.Error(err)
	}
	t.Log("get sql.DB success")

	conn1, err := db.Conn(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log("get one conn annot close")
	t.Log("begin get second conn")

	var begin time.Time
	go func() {
		begin = time.Now()
		<-time.After(5 * time.Second)
		conn1.Close()
	}()

	// wait here. wait for conn1 closed or wait for ctx canceled
	_, err = db.Conn(context.Background())
	if err == nil && time.Now().Before(begin.Add(5*time.Second)) {
		t.Log("get second conn success")
	} else {
		t.Error("get second conn success, expect failed")
	}
}
