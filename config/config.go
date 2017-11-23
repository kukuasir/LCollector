package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var SysConfig = TomlConfig{}

type TomlConfig struct {
	Server  server      `toml:"server"`
	Mysql   dbConfig    `toml:"mysql"`
	MongoDB dbConfig    `toml:"mongodb"`
	Logger  logConfig   `toml:"logger"`
	Token   tokenConfig `toml:"token"`
}

type server struct {
	Port string `toml:"port"`
}

type dbConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	DBName   string `toml:"name"`
	UserName string `toml:"username"`
	Password string `toml:"password"`
	ConnMax  int64  `toml:"connection_max"`
}

type logConfig struct {
	EnableOperateLog bool `toml:"enable_operate_log"`
	EnableMessageLog bool `toml:"enable_message_log"`
}

type tokenConfig struct {
	ValidTime int64 `toml:"valid_secs"`
}

func ConfigGet() {

	fmt.Println("Read Config...")

	var config TomlConfig
	if _, err := toml.DecodeFile("config/config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
}
