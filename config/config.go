package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var App appConfig
var System systemConfig
var Mongo mongoConfig
var Logger logConfig

type TomlConfig struct {
	AppConfig    appConfig    `toml:"app"`    // App信息配置
	SystemConfig systemConfig `toml:"system"` // 系统设置信息
	MongoConfig  mongoConfig  `toml:"mongo"`  // mongo配置信息
	LogConfig    logConfig    `toml:"logger"` // 日志设置信息
}

type appConfig struct {
	Name      string `toml:"name"`      // App名字
	Logo      string `toml:"logo"`      // 系统Logo
	Summary   string `toml:"summary"`   // 系统描述
	Version   string `toml:"version"`   // 系统版本
	Copyright string `toml:"copyright"` // 版权
	QQ        string `toml:"qq"`        // QQ
	Wechat    string `toml:"wechat"`    // 微信公众号
	Website   string `toml:"website"`   // 网站
}

type systemConfig struct {
	ServerPort string `toml:"server_port"` // API服务端口
	ValidSecs  int64  `toml:"valid_times"` // Token有效时长
}

type mongoConfig struct {
	Host     string `toml:"host"`
	Database string `toml:"database"`
	UserName string `toml:"username"`
	Password string `toml:"password"`
}

type logConfig struct {
	EnableOperateLog bool   `toml:"enable_operate_log"` // 是否打开操作日志
	EnableMessageLog bool   `toml:"enable_message_log"` // 是否打开消息日志
	LogPath          string `toml:"log_path"`           // 本地日志文件存储路径
}

func InitConfig() {

	fmt.Println("Read Config...")

	var tomlConfig TomlConfig
	if _, err := toml.DecodeFile("config/config.toml", &tomlConfig); err != nil {
		panic(err)
	}

	App = tomlConfig.AppConfig
	System = tomlConfig.SystemConfig
	Mongo = tomlConfig.MongoConfig
	Logger = tomlConfig.LogConfig
}
