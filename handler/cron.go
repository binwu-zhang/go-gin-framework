package handler

import (
	"github.com/robfig/cron"
)

func RegisterCron() {
	c := cron.New()
	//Field name   | Mandatory? | Allowed values  | Allowed special characters
	//----------   | ---------- | --------------  | --------------------------
	//Minutes      | Yes        | 0-59            | * / , -
	//Hours        | Yes        | 0-23            | * / , -
	//Day of month | Yes        | 1-31            | * / , - ?
	//Month        | Yes        | 1-12 or JAN-DEC | * / , -
	//Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

	//更新缓存 10秒执行一次
	//atomicUpdateCacheByDatabase := int32(0)
	//_ = c.AddFunc("*/10 * * * * ?", func() {
	//	if atomic.LoadInt32(&atomicUpdateCacheByDatabase) != 0 {
	//		return
	//	}
	//	atomic.StoreInt32(&atomicUpdateCacheByDatabase, 1)
	//	defer func() {
	//		atomic.StoreInt32(&atomicUpdateCacheByDatabase, 0)
	//	}()
	//	new(script.CacheDataScript).UpdateCacheByDatabase()
	//	return
	//})

	//每天凌晨3点执行 全量更新本地缓存
	//_ = c.AddFunc("0 0 3 * * ?", new(script.CacheDataScript).UpdateAllCache)

	c.Start()
}
