package db

import (
	"aig-tech-okr/libs"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var mysqldbs = make(map[string]*gorm.DB)
var rwm sync.RWMutex

func RegisterMysqlConnectPool() error {
	rwm.Lock()
	defer rwm.Unlock()

	for _, v := range libs.Conf.Mysql {
		dsn := mysql.Config{
			DSN: fmt.Sprintf("%s:%s@%s(%s)/%s?charset=%s&parseTime=True&loc=Local&allowNativePasswords=True", v.Username, v.Password, v.Net, v.Host, v.Dbname, v.Charset),
		}

		db, err := gorm.Open(mysql.New(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return errors.New("[" + v.InsName + "]open失败:" + err.Error())
		}
		if sqlDb, err := db.DB(); err != nil {
			sqlDb.SetMaxOpenConns(200)
			sqlDb.SetMaxIdleConns(100)
			sqlDb.SetConnMaxLifetime(time.Minute * 5)
			if err = sqlDb.Ping(); err != nil {
				return errors.New("[" + v.InsName + "]ping失败:" + err.Error())
			}
		}
		mysqldbs[v.InsName] = db
	}
	return nil
}

func connectdb(dbname string) *gorm.DB {
	logKey := "mysql-connectdb"
	for _, v := range libs.Conf.Mysql {
		if v.InsName == dbname {
			dsn := mysql.Config{
				DSN: fmt.Sprintf("%s:%s@%s(%s)/%s?charset=%s&parseTime=True&loc=Local&allowNativePasswords=True", v.Username, v.Password, v.Net, v.Host, v.Dbname, v.Charset),
			}

			db, err := gorm.Open(mysql.New(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			if err != nil {
				libs.ErrorLog(logKey, "["+v.InsName+"]open失败:", "", err)
				return &gorm.DB{}
			}
			if sqlDb, err := db.DB(); err != nil {
				sqlDb.SetMaxOpenConns(200)
				sqlDb.SetMaxIdleConns(100)
				sqlDb.SetConnMaxLifetime(time.Minute * 5)
				if err = sqlDb.Ping(); err != nil {
					libs.ErrorLog(logKey, "连接数据库失败", "", err)
					return &gorm.DB{}
				}
			}

			return db
		}
	}

	libs.ErrorLog(logKey, "mysql配置不存在："+dbname, "", nil)
	return &gorm.DB{}
}

func newConnect(dbName string) *gorm.DB {
	rwm.Lock()
	defer rwm.Unlock()

	if db, ok := mysqldbs[dbName]; ok && db.Config != nil {
		sqlDb, err := db.DB()
		if err == nil && sqlDb.Ping() == nil {
			return db
		}
	}

	db := connectdb(dbName)
	mysqldbs[dbName] = db
	return db
}

func getConnect(dbName string) *gorm.DB {
	rwm.RLock()
	defer rwm.RUnlock()
	return mysqldbs[dbName]
}

func GetMysqlConn(dbName string) *gorm.DB {

	db := getConnect(dbName)
	if db == nil || db.Config == nil {
		return newConnect(dbName)
	}

	sqlDb, err := db.DB()
	if err == nil && sqlDb.Ping() == nil {
		return db
	}
	return newConnect(dbName)
}
