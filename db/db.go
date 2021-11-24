package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
)

type DBPool struct {
	dbs map[string]*MyDB
}

type MyDB struct {
	*gorm.DB
	dbType  string
	key     string
	name    string
	options *Options
}

type Options struct {
	maxConn     int
	idleConn    int
	maxOpenConn int
	maxLeftTime time.Duration
}

type Option func(*Options)

func NewDBPool() *DBPool {
	return &DBPool{
		dbs: make(map[string]*MyDB),
	}
}

// AddDB add a db instance into poll
// note that, all instance should be added in init status of Application
//
// example:
// AddDB("user", "root:root@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local")
// then, you can use code `GetDB("user")` to get it
func (p DBPool) AddDB(dbType string, key string, dsn string, opts ...Option) error {
	options := &Options{
		maxConn:     100,
		idleConn:    100,
		maxLeftTime: 60 * time.Minute,
	}
	options.configure(opts...)

	//create db if not exist
	createDBIfNotExist(dbType, dsn)

	//reconnect and init instance
	gormInstance, err := connect(dbType, dsn)
	if err != nil {
		return err
	}

	db, err := gormInstance.DB()
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(options.maxConn)
	db.SetMaxIdleConns(options.idleConn)
	db.SetMaxOpenConns(options.maxOpenConn)

	// Mysql usually close a conn if it doesn't use in eight hours,
	// So maxLeftTime best less than eight hours,
	// We close it active if don't use in maxLeftTime. avoid it turn into a broke pipe
	db.SetConnMaxLifetime(options.maxLeftTime)

	if err = db.Ping(); err != nil {
		return err
	}

	i := strings.LastIndex(dsn, "/")
	dbName := dsn[:1+i]
	p.dbs[key] = &MyDB{
		DB:      gormInstance,
		dbType:  dbType,
		key:     key,
		name:    dbName,
		options: options,
	}

	return nil
}

func connect(dbType string, dsn string) (*gorm.DB, error) {
	switch dbType {
	case MYSQL:
		return gorm.Open(mysql.New(mysql.Config{
			DSN: dsn,
		}), &gorm.Config{})
	case POSTGRES:
		return gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
		}), &gorm.Config{Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 500 * time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      true,
		})})
	default:
		return nil, errors.New("not implement db type")
	}
}

func createDBIfNotExist(dbType string, dsn string) error {
	// postgres need not create db
	if dbType == POSTGRES {
		return nil
	}

	i := strings.LastIndex(dsn, "/")
	db, err := connect(dbType, dsn[:i+1])
	if err != nil {
		return err
	}

	j := strings.Index(dsn, "?")
	if j == -1 {
		j = len(dsn)
	}

	sql := fmt.Sprintf("create database if not exists %s default character set utf8mb4 collate utf8mb4_unicode_ci;", dsn[i+1:j])
	if err = db.Exec(sql).Error; err != nil {
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return err
	}
	if err = sqlDb.Close(); err != nil {
		return err
	}
	return nil
}

func (p DBPool) GetDB(name string) *MyDB {
	if db, ok := p.dbs[name]; ok {
		return db
	}

	return nil
}

func (p DBPool) ReleasePool() {
	for _, instance := range p.dbs {
		db, _ := instance.DB.DB()
		_ = db.Close()
	}
}

func (o *Options) configure(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// --------------------------------------------------------------------

var defaultPool *DBPool

func init() {
	defaultPool = NewDBPool()
}

func AddDB(dbType string, key string, sqlInfo string, opts ...Option) error {
	return defaultPool.AddDB(dbType, key, sqlInfo, opts...)
}

func GetDB(key string) *MyDB {
	return defaultPool.GetDB(key)
}

func MustGetDB(key string) *MyDB {
	db := GetDB(key)
	if db == nil {
		panic("MyDB " + key + " not exist")
	}
	return db
}

func ReleaseDBPool() {
	defaultPool.ReleasePool()
}

func MaxConn(n int) Option {
	return func(o *Options) {
		o.maxConn = n
	}
}

func IdleConn(n int) Option {
	return func(o *Options) {
		o.idleConn = n
	}
}

func MaxOpenConn(n int) Option {
	return func(o *Options) {
		o.maxOpenConn = n
	}
}

func MaxLeftTime(duration int64) Option {
	return func(o *Options) {
		o.maxLeftTime = time.Duration(duration) * time.Minute
	}
}
