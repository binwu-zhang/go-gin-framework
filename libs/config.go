package libs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

// Conf 全局配置
var Conf config

type config struct {
	App   appConf     `mapstructure:"app"`
	Log   logConf     `mapstructure:"log"`
	Mysql []mysqlConf `mapstructure:"mysql"`
	Redis []redisConf `mapstructure:"redis"`
	Mongo []mongoConf `mapstructure:"mongo"`
}

// mongo数据库配置
type mongoConf struct {
	InsName     string `mapstructure:"ins_name"`
	Addr        string `mapstructure:"addr"`
	PemPath     string `mapstructure:"pem_path"`
	MaxPoolSize int    `mapstructure:"max_pool_size"`
}

// 日志配置
type logConf struct {
	LogPath   string `mapstructure:"log_path"`
	ReqPath   string `mapstructure:"req_path"`
	CronPath  string `mapstructure:"cron_path"`
	InfoPath  string `mapstructure:"info_path"`
	ErrorPath string `mapstructure:"error_path"`
}

// app配置
type appConf struct {
	SignPrefix   string `mapstructure:"sign_prefix"`
	SignSuffix   string `mapstructure:"sign_suffix"`
	ServerAddr   string `mapstructure:"server_addr"`
	CronHostName string `mapstructure:"cron_host_name"`
	Env          string `mapstructure:"env"`
}

// mysql数据库配置
type mysqlConf struct {
	InsName  string `mapstructure:"ins_name"`
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	Net      string `mapstructure:"net"`
	Charset  string `mapstructure:"charset"`
}

// redis配置
type redisConf struct {
	InsName      string `mapstructure:"ins_name"`
	Addr         string `mapstructure:"addr"`
	Auth         string `mapstructure:"auth"`
	Db           int    `mapstructure:"db"`
	ConnTimeout  int    `mapstructure:"conn_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	MaxIdle      int    `mapstructure:"max_idle"`
	MaxActive    int    `mapstructure:"max_active"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
	MaxWait      bool   `mapstructure:"max_wait"`
}

func RegisterConfig() {

	//读取yaml文件
	viper.AddConfigPath("./config")

	env := os.Getenv("BUSINESS_ENV")
	fmt.Println(env)
	if env == "" {
		env = "product"
	}

	viper.SetConfigName(env)

	//设置配置文件类型
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic("初始化config失败-读配置:" + err.Error())
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic("初始化config失败-解析:" + err.Error())
	}

	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		if err := viper.Unmarshal(&Conf); err != nil {
			fmt.Println("配置变更失败-解析:" + err.Error())
		}
	})
}
