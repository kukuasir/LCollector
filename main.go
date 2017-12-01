package main

import (
	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//
	"LCollector/config"
	sw "LCollector/controller"
	"log"
	"net/http"
)

func main() {

	log.Printf("Server started")

	// 读取配置文件
	config.InitConfig()

	// 自定义错误码
	config.InitErrors()

	// 创建路由
	router := sw.NewRouter()

	err := http.ListenAndServe(config.System.ServerPort, router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
