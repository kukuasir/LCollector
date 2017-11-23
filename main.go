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
	"strings"
)

func main() {

	log.Printf("Server started")

	// 读取配置文件
	config.ConfigGet()

	router := sw.NewRouter()

	port := strings.Join([]string{":"}, config.SysConfig.Server.Port)
	log.Fatal(http.ListenAndServe(port, router))
}
