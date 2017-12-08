package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/** 定义操作类型 */
const (
	OPERATE_TYPE_ADD    = 1
	OPERATE_TYPE_DELETE = 2
	OPERATE_TYPE_UPDATE = 3
)

/** 操作对象 */
const (
	OPERATE_TARGET_USER     = 1
	OPERATE_TARGET_AGENCY   = 2
	OPERATE_TARGET_DEVICE   = 3
	OPERATE_TARGET_PASSWORD = 4
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

// 查询列表的总个数
func GetCount(coll string) (int64, error) {
	var count int
	query := func(c *mgo.Collection) error {
		var err error
		count, err = c.Find(bson.M{}).Count()
		return err
	}
	err := SharedQuery(coll, query)
	return int64(count), err
}

// 校验page的值
func ValidPageValue(page *int) {
	if *page > 0 {
		*page = *page - 1
	} else {
		*page = 0
	}
}
