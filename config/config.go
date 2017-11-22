package config

import (
	"github.com/larspensjo/config"
	"fmt"
)

var ServerPort string

var EnableOperateLog bool
var EnableMessageLog bool

var TokenValidTime int

func ReadSystemConfig() {

	fmt.Println("Read Config...")

	c, _ := config.ReadDefault("config.cfg")

	ServerPort, _ = c.String("SERVER", "SERVER_PORT")

	host, _ := c.String("MONGODB", "DB_HOST")
	port, _ := c.String("MONGODB", "DB_PORT")
	name, _ := c.String("MONGODB", "DB_NAME")
	username, _ := c.String("MONGODB", "DB_USERNAME")
	//password, _ := c.String("MONGODB", "DB_PASSWORD")
	fmt.Println(username + "@" + host + ":" + port + "/" + name)

	EnableOperateLog, _ = c.Bool("LOG", "ENABLE_OPERATE_LOG")
	EnableMessageLog, _ = c.Bool("LOG", "ENABLE_MESSAGE_LOG")

	TokenValidTime, _ = c.Int("TOKEN", "TOKEN_VALID_TIME")
}
