package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

/** 定义操作类型 */
const (
	OPERATE_TYPE_ADD    = 1
	OPERATE_TYPE_DELETE = 2
	OPERATE_TYPE_UPDATE = 3
)

/** 操作对象 */
const (
	OPERATE_TARGET_USER   = 10
	OPERATE_TARGET_AGENCY = 11
	OPERATE_TARGET_DEVICE = 12
)

func WriteData(w http.ResponseWriter, res interface{}) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Fatal("json marshal error: ", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Server", "BTK")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	defer func() {
		w.Write(data)
	}()
}
