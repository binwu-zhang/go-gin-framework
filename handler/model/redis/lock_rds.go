package redis

import (
	"aig-tech-okr/libs"
	"aig-tech-okr/libs/cache"
	"fmt"
	"time"
)

const (
	defaultRetryIntervalTime = 500
	defaultRetryTimes        = 1
)

//LockRds
//  @Description: 通用锁
type LockRds struct {
	traceId           string
	key               string
	expireTime        time.Duration //锁过期时间，过期自动删除，单位是秒 =0永久不删除 >0则超时删除
	isRetry           bool          //是否重试 true:是 false:否 默认:否,不重试
	retryIntervalTime int           //加锁重试间隔时间 单位:毫秒 默认:500
	retryTimes        int           //加锁重试次数 默认:1
}

func (i *LockRds) ins() string {
	return "okr"
}

func NewLock(traceId string) *LockRds {
	return &LockRds{
		traceId: traceId,
	}
}

//redis锁 isSuccess true:加锁成功 false:加锁失败
func (i *LockRds) lock() (isSuccess bool, err error) {

	if i.isRetry {
		return i.retryLock()
	}

	rds := cache.PoolGet(i.ins())

	res, err := rds.SetNX(i.key, "1", i.expireTime).Result()

	if err == nil { //加锁成功

		return res, nil

	} else { //redis报错

		libs.ErrorLog("lockRds-lock", "加锁报错", i.traceId, err)

		return false, err
	}

}

func (i *LockRds) retryLock() (isSuccess bool, err error) {
	rds := cache.PoolGet(i.ins())

	if i.retryIntervalTime == 0 {
		i.retryIntervalTime = defaultRetryIntervalTime
	}
	if i.retryTimes == 0 {
		i.retryIntervalTime = defaultRetryTimes
	}

	for times := 0; times <= i.retryTimes; times++ {

		err := rds.SetNX(i.key, "1", i.expireTime).Err()

		if err == nil { //加锁成功

			return true, nil

		} else { //redis报错

			libs.ErrorLog("lockRds-lock", "加锁报错", i.traceId, err)

			return false, err
		}
	}

	return
}

func (i *LockRds) Unlock() error {

	rds := cache.PoolGet(i.ins())

	if err := rds.Del(i.key).Err(); err != nil {

		libs.ErrorLog("lockRds-unlock", "删锁报错", i.traceId, err)

		return err
	}

	return nil
}

////锁，防止并发
//	lock := redis.NewLock(i.TraceId)
//	if lockSuccess, err := lock.user(id); !lockSuccess {
//		libs.ErrorLog(logKey, "加锁失败", i.TraceId, err)
//		return
//	}
//	defer func() {
//		_ = lock.Unlock()
//		if err := recover(); err != nil {
//			libs.ErrorLog(logKey, "panic", i.TraceId, err)
//		}
//	}()

func (i *LockRds) user(id int) (isSuccess bool, err error) {
	i.expireTime = 3 * time.Second
	i.key = fmt.Sprintf("user:lock:%d", id)
	return i.lock()
}
