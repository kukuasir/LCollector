package controller

import (
	"LCollector/config"
	"LCollector/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
	"strings"
)

func FetchMessageLogList(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size == 0 {
		size = 20 // 默认一页加载20条数据
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}

	// 只有超级管理员才有权限查看消息日志
	if operator.Role != "root" {
		WriteData(w, config.NewError(config.PermissionDeniedAgency))
		return
	}

	logList, err := fetchPagingMessageLogs(operator, page, size)

	// 计算数据总条数
	var totalCount int64
	if page == 0 {
		totalCount, err = GetCount(T_MESSAGE_LOG)
	}

	// 返回查询结果
	var logListRet model.MessageLogRet
	logListRet.ResultInfo.Status = config.Success
	logListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	logListRet.ResultInfo.Total = totalCount
	logListRet.MessageList = logList
	WriteData(w, logListRet)
}

func FetchOperateLogList(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size == 0 {
		size = 20 // 默认一页加载20条数据
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}

	tempLogs, err := fetchPagingOperateLogs(operator, page, size)
	if err != nil {
		panic(err)
	}

	// 转换到操作日志表中
	var operateLogs []model.OperateLog
	for i := 0; i < len(tempLogs); i++ {
		temp := tempLogs[i]
		var log model.OperateLog
		log.Type = temp.Type
		log.Target = temp.Target
		log.TargetObject = temp.TargetObject
		log.OperatorId = temp.OperatorId
		log.CreateTime = temp.CreateTime
		log.SourceIP = temp.SourceIP
		log.AgencyId = temp.AgencyId
		if len(temp.UserNames) > 0 {
			log.OperatorName = temp.UserNames[0]
		}
		operateLogs = append(operateLogs, log)
	}

	// 计算数据总条数
	var totalCount int64
	if page == 0 {
		totalCount, err = GetCount(T_OPERATE_LOG)
	}

	// 返回查询结果
	var logListRet model.OperateLogRet
	logListRet.ResultInfo.Status = config.Success
	logListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	logListRet.ResultInfo.Total = totalCount
	logListRet.OperateList = operateLogs
	WriteData(w, logListRet)
}

func FetchLoginLogList(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size == 0 {
		size = 20 // 默认一页加载20条数据
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}

	tempLogs, err := fetchPagingLoginLogs(operator, page, size)
	if err != nil {
		panic(err)
	}

	// 转换到登录日志表中
	var loginLogs []model.LoginLog
	for i := 0; i < len(tempLogs); i++ {
		temp := tempLogs[i]
		var log model.LoginLog
		log.UserId = temp.UserId
		log.StatusDesc = config.UserStatusDesc(temp.Status)
		log.AgencyId = temp.AgencyId
		log.CreateTime = temp.CreateTime
		log.SourceIP = temp.SourceIP
		if len(temp.UserNames) > 0 {
			log.UserName = temp.UserNames[0]
		}
		loginLogs = append(loginLogs, log)
	}

	// 计算数据总条数
	var totalCount int64
	if page == 0 {
		totalCount, err = GetCount(T_LOGIN_LOG)
	}

	// 返回查询结果
	var logListRet model.LoginLogRet
	logListRet.ResultInfo.Status = config.Success
	logListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	logListRet.ResultInfo.Total = totalCount
	logListRet.LoginList = loginLogs
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
func InsertMessageLog(msgType int64, deviceNo string, content string, ipaddr string) bool {
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"type":        msgType,
			"type_desc":   msgTypeDesc(msgType),
			"device_no":   deviceNo,
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

func fetchPagingLoginLogs(operator model.User, page, size int) ([]model.TempLoginLog, error) {

	var tempLogs []model.TempLoginLog

	if operator.Role == "customer" {
		return tempLogs, nil
	}

	ValidPageValue(&page)

	query := func(c *mgo.Collection) error {
		pipeline := []bson.M{
			bson.M{"$sort": bson.M{"create_time": -1}},
			bson.M{"$skip": page * size},
			bson.M{"$limit": size},
			bson.M{"$lookup": bson.M{"from": T_USER, "localField": "user_id", "foreignField": "_id", "as": "user_docs"}},
			bson.M{"$project": bson.M{
				"user_id":     1,
				"status":      1,
				"agency_id":   1,
				"create_time": 1,
				"source_ip":   1,
				"user_names":  "$user_docs.user_name",
			}},
		}
		if operator.Role == "admin" {
			pipeline = append(pipeline, bson.M{"$match": bson.M{"agency_id": operator.AgencyId}})
		}
		return c.Pipe(pipeline).All(&tempLogs)
	}
	err := SharedQuery(T_LOGIN_LOG, query)
	return tempLogs, err
}

func fetchPagingOperateLogs(operator model.User, page, size int) ([]model.TempOperateLog, error) {

	var tempLogs []model.TempOperateLog

	if operator.Role == "customer" {
		return tempLogs, nil
	}

	ValidPageValue(&page)

	query := func(c *mgo.Collection) error {
		pipeline := []bson.M{
			bson.M{"$sort": bson.M{"create_time": -1}},
			bson.M{"$skip": page * size},
			bson.M{"$limit": size},
			bson.M{"$lookup": bson.M{"from": T_USER, "localField": "operator_id", "foreignField": "_id", "as": "user_docs"}},
			bson.M{"$project": bson.M{
				"type":        1,
				"target":      1,
				"object":      1,
				"operator_id": 1,
				"agency_id":   1,
				"create_time": 1,
				"source_ip":   1,
				"user_names":  "$user_docs.user_name",
			}},
		}
		if operator.Role == "admin"  {
			pipeline = append(pipeline, bson.M{"$match": bson.M{"agency_id": operator.AgencyId}})
		}
		return c.Pipe(pipeline).All(&tempLogs)
	}
	err := SharedQuery(T_OPERATE_LOG, query)
	return tempLogs, err
}

func fetchPagingMessageLogs(operator model.User, page, size int) ([]model.MessageLog, error) {

	var loglist []model.MessageLog

	if operator.Role != "root" {
		return loglist, nil
	}

	ValidPageValue(&page)

	query := func(c *mgo.Collection) error {
		return c.Find(nil).Sort("-create_time").Skip(page * size).Limit(size).All(&loglist)
	}
	err := SharedQuery(T_MESSAGE_LOG, query)
	return loglist, err
}

func msgTypeDesc(msgtype int64) string {
	if msgtype == model.MESSAGE_TYPE_HEARTBEAT {
		return "HEARTBEAT"
	} else if msgtype == model.MESSAGE_TYPE_STATUS {
		return "STATUS"
	} else if msgtype == model.MESSAGE_TYPE_CONFIG {
		return "CONFIG"
	} else if msgtype == model.MESSAGE_TYPE_DATA {
		return "DATA"
	} else if msgtype == model.MESSAGE_TYPE_WARNING {
		return "WARNING"
	} else if msgtype == model.MESSAGE_TYPE_REQUEST {
		return "REQUEST"
	}
	return "NONE"
}
