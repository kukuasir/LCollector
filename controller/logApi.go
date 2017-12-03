package controller

import (
	"LCollector/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func FetchMessageLogList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func FetchOperateLogList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func FetchLoginLogList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// 插入登录日志
func InsertLoginLog(user model.User, ipaddr string) bool {
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"user_id":   user.UserId.Hex(),
			"status":    user.Status,
			"agency_id": user.AgencyId,
			"time":      time.Now().Unix(),
			"on_ip":     ipaddr,
		}
		return c.Insert(selector)
	}
	err := SharedQuery(T_LOGIN_LOG, query)
	if err != nil {
		return false
	}
	return true
}

// 插入操作日志
func InsertOperateLog(operateType int64, userId string, agencyId string, target string, ipaddr string) bool {
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"operator_type": operateType,
			"operator_id":   userId,
			"agency_id":     agencyId,
			"target":        target,
			"time":          time.Now().Unix(),
			"onip":          ipaddr,
		}
		return c.Insert(selector)
	}
	err := SharedQuery(T_OPERATE_LOG, query)
	if err != nil {
		return false
	}
	return true
}

// 插入消息日志
func InsertMessageLog(deviceId string, agencyId string, content string, ipaddr string) bool {
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"device_id": deviceId,
			"agency_id": agencyId,
			"content":   content,
			"time":      time.Now().Unix(),
			"onip":      ipaddr,
		}
		return c.Insert(selector)
	}
	err := SharedQuery(T_MESSAGE_LOG, query)
	if err != nil {
		return false
	}
	return true
}
