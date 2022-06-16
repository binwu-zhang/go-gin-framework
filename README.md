

### 目录：

- config：配置文件
- handler：
    - controller：接收请求参数，处理相关逻辑
    - action：常用方法
    - middleware：中间件
    - entity：结构体定义
    - cont：常量定义
    - model：数据库model
    - service：调用第三方服务
    - script：脚本
    - util：公用方法，和业务无耦合
    - cron.go：定时任务路由
    - router.go：API路由
    - script.go：脚本路由
- libs：自定义库

### 日志：

- 错误日志：libs.ErrorLog(logKey, desc, traceId, err)
- 调试日志：libs.InfoLog(logKey, desc, traceId, err)
- 脚本日志：libs.CronLog(logKey, desc, traceId, err)

### 三方扩展：
- [gin](github.com/gin-gonic/gin)
- [gormV2](https://gorm.io/zh_CN/)
    - [文档](https://gorm.io/zh_CN/docs/)
    - [v1和v2区别](https://gorm.io/zh_CN/docs/changelog.html)
- [freecache](github.com/coocood/freecache)
- [定时任务](github.com/robfig/cron)