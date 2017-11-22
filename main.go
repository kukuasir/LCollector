package main

import (
	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//
	sw "./controller"
	"log"
	"net/http"
	"LCollector/config"
)

func main() {

	log.Printf("Server started")

	// 读取配置文件
	config.ReadSystemConfig()

	router := sw.NewRouter()
	
	log.Fatal(http.ListenAndServe(":" + config.ServerPort, router))
}
