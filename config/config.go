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
	AppConfig    appConfig    `toml:"app"`
	SystemConfig systemConfig `toml:"system"`
	MongoConfig  mongoConfig  `toml:"mongo"`
	LogConfig    logConfig    `toml:"logger"`
}

type appConfig struct {
	Name      string `toml:"name"`
	Logo      string `toml:"logo"`
	Version   string `toml:"version"`
	Copyright string `toml:"copyright"`
	QQ        string `toml:"qq"`
	Wechat    string `toml:"wechat"`
	Website   string `toml:"website"`
}

type systemConfig struct {
	ServerPort string `toml:"server_port"`
	ValidSecs  int64  `toml:"valid_times"`
}

type mongoConfig struct {
	Host     string `toml:"host"`
	Database string `toml:"database"`
	UserName string `toml:"username"`
	Password string `toml:"password"`
}

type logConfig struct {
	EnableOperateLog bool   `toml:"enable_operate_log"`
	EnableMessageLog bool   `toml:"enable_message_log"`
	LogPath          string `toml:"log_path"`
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
