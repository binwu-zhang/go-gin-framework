package main

import (
	"aig-tech-okr/handler"
	"aig-tech-okr/handler/action"
	"aig-tech-okr/handler/middleware"
	"aig-tech-okr/handler/util"
	"aig-tech-okr/libs"
	"aig-tech-okr/libs/cache"
	"aig-tech-okr/libs/db"
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"time"
)

func main() {

	gin.DisableConsoleColor()

	//注册配置文件
	libs.RegisterConfig()

	//注册log
	libs.RegisterLog()

	//c := runtime.GOMAXPROCS(2)
	//libs.InfoLog("main", fmt.Sprintf("%d", c), "", nil)

	var hostName, _ = os.Hostname()
	var env = os.Getenv("AIGOA_BUSINESS_ENV")
	var traceId = hostName + "_" + util.UniqueID() + "_" + env + "_" + strconv.Itoa(os.Getpid())

	defer func() {
		if r := recover(); r != nil {
			buf := debug.Stack()
			libs.ErrorLog("main", "服务错误", traceId, fmt.Sprintln(fmt.Sprintf("%s", buf)))
			fmt.Println("服务panic退出", traceId, r)
			return
		}
		libs.InfoLog("main", "服务关闭", traceId, nil)
	}()

	//注册redis
	if err := cache.RegisterRedisPool(); err != nil {
		fmt.Println("main_server", "服务启动失败-注册redisPool", traceId, err)
		libs.ErrorLog("main_server", "服务启动失败-注册redisPool", traceId, err)
		return
	}

	db.RegisterMongo()

	//注册mysql
	if err := db.RegisterMysqlConnectPool(); err != nil {
		fmt.Println("main_server", "服务启动失败-注册mysqlPool", traceId, err)
		libs.ErrorLog("main_server", "服务启动失败-注册mysqlPool", traceId, err)
		return
	}
	//初始化缓存
	action.InitCache()

	//脚本入口
	if len(os.Args) > 1 {
		handler.RegisterScript()
		return
	}

	//初始化RBAC
	//libs.InitRBAC()

	//注册cron
	handler.RegisterCron()

	//gin.SetMode(gin.DebugMode)
	//engine := gin.New()
	engine := gin.Default()

	//跨域
	engine.Use(middleware.Cors())
	//engine.Use(middleware.Metric())

	//请求日志记录
	engine.Use(middleware.Logger())

	//异常处理
	engine.Use(middleware.CustomRecovery())
	//请求速率限制
	//engine.Use(middleware.MaxAllowed(50))

	//注册路由
	handler.RegisterRouter(engine)
	handler.RegisterRouterV2(engine)

	//pprof
	pprof.Register(engine, "pprof")

	/*****************设置graceful shutdown*********************/
	srv := &http.Server{
		Addr:        libs.Conf.App.ServerAddr,
		Handler:     engine,
		ReadTimeout: time.Second * 3,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
