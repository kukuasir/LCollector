package controller

import (
	"LCollector/config"
	"LCollector/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FetchMessageLogList(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	// 验证操作人是否存在
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	logList, err := fetchPagingMessageLogs(operator, page, size)
	if err != nil {
		panic(err)
	}

	// 返回查询结果
	var logListRet model.MessageLogRet
	logListRet.ResultInfo.Status = config.Success
	logListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	logListRet.MessageList = logList
	WriteData(w, logListRet)
}

func FetchOperateLogList(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	// 验证操作人是否存在
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	logList, err := fetchPagingOperateLogs(operator, page, size)
	if err != nil {
		panic(err)
	}

	// 返回查询结果
	var logListRet model.OperateLogRet
	logListRet.ResultInfo.Status = config.Success
	logListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	logListRet.OperateList = logList
	WriteData(w, logListRet)
}

func FetchLoginLogList(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	// 验证操作人是否存在
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	logList, err := fetchPagingLoginLogs(operator, page, size)
	if err != nil {
		panic(err)
	}

	// 返回查询结果
	var logListRet model.LoginLogRet
	logListRet.ResultInfo.Status = config.Success
	logListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	logListRet.LoginList = logList
	WriteData(w, logListRet)
}

// 插入登录日志
func InsertLoginLog(user model.User, ipaddr string) bool {
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"user_id":     user.UserId,
			"status":      user.Status,
			"agency_id":   user.AgencyId,
			"create_time": time.Now().Unix(),
			"source_ip":   ipaddr,
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
func InsertOperateLog(mode int64, target int64, operator model.User, object string, ipaddr string) bool {
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"type":        mode,
			"target":      target,
			"operator_id": operator.UserId,
			"agency_id":   operator.AgencyId,
			"object":      object,
			"create_time": time.Now().Unix(),
			"source_ip":   ipaddr,
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
			"device_id":   deviceId,
			"agency_id":   agencyId,
			"content":     content,
			"create_time": time.Now().Unix(),
			"source_ip":   ipaddr,
		}
		return c.Insert(selector)
	}
	err := SharedQuery(T_MESSAGE_LOG, query)
	if err != nil {
		return false
	}
	return true
}

////=========== Private Methods ===========

func fetchPagingLoginLogs(operator model.User, page, size int) ([]model.LoginLog, error) {

	var loglist []model.LoginLog

	if operator.Role == "customer" {
		return nil, nil
	}

	resp := bson.M{}
	query := func(c *mgo.Collection) error {
		var pipeline []bson.M
		if operator.Role == "root" {
			pipeline = []bson.M{
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"$unwind": "$agency"},
				bson.M{"$sort": bson.M{"time": -1}},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		} else if operator.Role == "admin" {
			pipeline = []bson.M{
				bson.M{"agency_id": operator.AgencyId},
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"$unwind": "$agency"},
				bson.M{"$sort": bson.M{"time": -1}},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		}
		return c.Pipe(pipeline).All(&resp)
	}
	err := SharedQuery(T_LOGIN_LOG, query)
	return loglist, err
}

func fetchPagingOperateLogs(operator model.User, page, size int) ([]model.OperateLog, error) {

	var loglist []model.OperateLog

	if operator.Role == "customer" {
		return nil, nil
	}

	resp := bson.M{}
	query := func(c *mgo.Collection) error {
		var pipeline []bson.M
		if operator.Role == "root" {
			pipeline = []bson.M{
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"unwind": "$agency"},
				bson.M{"$sort": bson.M{"time": -1}},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		} else if operator.Role == "admin" {
			pipeline = []bson.M{
				bson.M{"agency_id": operator.AgencyId},
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"unwind": "$agency"},
				bson.M{"$sort": bson.M{"time": -1}},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		}
		return c.Pipe(pipeline).All(&resp)
	}
	err := SharedQuery(T_OPERATE_LOG, query)
	return loglist, err
}

func fetchPagingMessageLogs(operator model.User, page, size int) ([]model.MessageLog, error) {

	var loglist []model.MessageLog

	if operator.Role == "customer" {
		return nil, nil
	}

	resp := bson.M{}
	query := func(c *mgo.Collection) error {
		var pipeline []bson.M
		if operator.Role == "root" {
			pipeline = []bson.M{
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"unwind": "$agency"},
				bson.M{"$sort": bson.M{"time": -1}},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		} else if operator.Role == "admin" {
			pipeline = []bson.M{
				bson.M{"agency_id": operator.AgencyId},
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"unwind": "$agency"},
				bson.M{"$sort": bson.M{"time": -1}},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		}
		return c.Pipe(pipeline).All(&resp)
	}
	err := SharedQuery(T_MESSAGE_LOG, query)
	return loglist, err
}
