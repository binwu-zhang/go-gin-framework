package cache

import (
	"aig-tech-okr/libs"
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"sync"
	"time"
)

var redisPool = make(map[string]*redis.Client)
var rwm sync.RWMutex

var redisClusterPool = make(map[string]*redis.ClusterClient)

func RegisterRedisPool() error {
	rwm.Lock()
	defer rwm.Unlock()

	redisConfs := libs.Conf.Redis

	for k, conf := range redisConfs {

		if conf.InsName == "" || conf.Addr == "" {
			return errors.New("[" + strconv.Itoa(k) + "]" + "实例名、地址不能为空")
		}

		if _, ok := redisPool[conf.InsName]; ok {
			continue
		}

		options := &redis.Options{
			Network:      "tcp",
			Addr:         conf.Addr,
			Password:     conf.Auth,
			DialTimeout:  time.Duration(conf.ConnTimeout) * time.Millisecond,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			PoolSize:     conf.MaxActive,
			MinIdleConns: 20,
			IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		}

		redisClient := redis.NewClient(options)

		_, err := redisClient.Ping().Result()
		if err != nil {
			return errors.New("[" + conf.InsName + "]ping失败:" + err.Error())
		}

		redisPool[conf.InsName] = redisClient
	}

	return nil
}

func connectRds(insName string) *redis.Client {

	redisConfs := libs.Conf.Redis

	for _, conf := range redisConfs {

		if conf.InsName != insName {
			continue
		}
		if _, ok := redisPool[conf.InsName]; ok {
			continue
		}

		options := &redis.Options{
			Network:      "tcp",
			Addr:         conf.Addr,
			Password:     conf.Auth,
			DialTimeout:  time.Duration(conf.ConnTimeout) * time.Millisecond,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			PoolSize:     conf.MaxActive,
			MinIdleConns: 1,
			IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		}

		redisClient := redis.NewClient(options)

		_, err := redisClient.Ping().Result()
		if err != nil {
			libs.ErrorLog("", "["+conf.InsName+"]ping失败:", "", err)
			return &redis.Client{}
		}

		return redisClient
	}

	return &redis.Client{}
}

func newConnect(insName string) *redis.Client {
	rwm.Lock()
	defer rwm.Unlock()

	if db, ok := redisPool[insName]; ok {
		if _, err := db.Ping().Result(); err == nil {
			return db
		}
	}

	db := connectRds(insName)
	redisPool[insName] = db
	return db
}

func get(insName string) *redis.Client {
	rwm.RLock()
	defer rwm.RUnlock()
	return redisPool[insName]
}

func PoolGet(insName string) *redis.Client {

	db := get(insName)
	if db == nil {
		return newConnect(insName)
	}
	if _, err := db.Ping().Result(); err == nil {
		return db
	}
	return newConnect(insName)
}


func RegisterRedisClusterPool() error {
	rwm.Lock()
	defer rwm.Unlock()

	redisConfs := libs.Conf.Redis

	for k, conf := range redisConfs {

		if conf.InsName == "" || conf.Addr == "" {
			return errors.New("[" + strconv.Itoa(k) + "]" + "实例名、地址不能为空")
		}

		if _, ok := redisClusterPool[conf.InsName]; ok {
			continue
		}

		options := &redis.ClusterOptions{
			Addrs:        strings.Split(conf.Addr, ","),
			Password:     conf.Auth,
			DialTimeout:  time.Duration(conf.ConnTimeout) * time.Millisecond,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			PoolSize:     conf.MaxActive,
			MinIdleConns: 20,
			IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		}

		redisClient := redis.NewClusterClient(options)

		_, err := redisClient.Ping().Result()
		if err != nil {
			return errors.New("[" + conf.InsName + "]ping失败:" + err.Error())
		}

		redisClusterPool[conf.InsName] = redisClient
	}

	return nil
}

func connectRdsCluster(insName string) *redis.ClusterClient {

	redisConfs := libs.Conf.Redis

	for _, conf := range redisConfs {

		if conf.InsName != insName {
			continue
		}
		if _, ok := redisClusterPool[conf.InsName]; ok {
			continue
		}

		options := &redis.ClusterOptions{
			Addrs:        strings.Split(conf.Addr, ","),
			Password:     conf.Auth,
			DialTimeout:  time.Duration(conf.ConnTimeout) * time.Millisecond,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			PoolSize:     conf.MaxActive,
			MinIdleConns: 20,
			IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		}

		redisClient := redis.NewClusterClient(options)

		_, err := redisClient.Ping().Result()
		if err != nil {
			libs.ErrorLog("", "["+conf.InsName+"]ping失败:", "", err)
			return &redis.ClusterClient{}
		}

		return redisClient
	}

	return &redis.ClusterClient{}
}


func newConnectCluster(insName string) *redis.ClusterClient {
	rwm.Lock()
	defer rwm.Unlock()

	if db, ok := redisClusterPool[insName]; ok {
		if _, err := db.Ping().Result(); err == nil {
			return db
		}
	}

	db := connectRdsCluster(insName)
	redisClusterPool[insName] = db
	return db
}

func getCluster(insName string) *redis.ClusterClient {
	rwm.RLock()
	defer rwm.RUnlock()
	return redisClusterPool[insName]
}

//func PoolGet(insName string) *redis.ClusterClient {
//
//	db := getCluster(insName)
//	if db == nil {
//		return newConnectCluster(insName)
//	}
//	if _, err := db.Ping().Result(); err == nil {
//		return db
//	}
//	return newConnectCluster(insName)
//}
