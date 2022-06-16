package libs

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

var logger *zap.SugaredLogger
var ReqLog *logrus.Logger

func ErrorLog(key, desc, traceId string, err interface{}) {
	logger.Errorw(key, "desc", desc, "trace_id", traceId, "err", err)
}

func CronLog(key, desc, traceId string, err interface{}) {
	logger.Infow(key, "desc", desc, "trace_id", traceId, "cron", err)
}

func InfoLog(key, desc, traceId string, detail interface{}) {
	logger.Infow(key, "desc", desc, "trace_id", traceId, "detail", detail)
}

func CustomLog(key string, keysAndValues ...interface{}) {
	logger.Infow(key, keysAndValues...)
}

func registerLogrus(logPath string) *logrus.Logger {

	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.Out = getWriter(logPath)
	return log
}

func RegisterLog() {

	//创建日志目录
	if err := os.MkdirAll(Conf.Log.LogPath, os.ModePerm); err != nil {
		fmt.Println("创建日志文件夹失败", fmt.Sprintf("%s", err))
	}

	ReqLog = registerLogrus(Conf.Log.ReqPath)
	tmEncode := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "modules",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     tmEncode,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(Conf.Log.InfoPath)
	errorWriter := getWriter(Conf.Log.ErrorPath)
	cronWriter := getWriter(Conf.Log.CronPath)

	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	cronLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(errorWriter), errorLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(cronWriter), cronLevel),
	)

	logger = zap.New(core, zap.AddCaller()).WithOptions(zap.AddCallerSkip(1)).Sugar()
}

//日志文件切割
func getWriter(filename string) io.Writer {
	// 保存30天内的日志，每24小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
