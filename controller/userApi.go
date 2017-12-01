package controller

import (
	"LCollector/config"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"LCollector/model"
	"LCollector/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func AddUser(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 &&
		strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.UserReq
	json.Unmarshal(body, &req)

	// 验证数据是否有效
	if !util.ValidAddUser(req) {
		WriteData(w, config.InvalidParameterValue)
		return
	}

	// 验证操作人信息是否有效
	operator, err := queryUserByID(req.OperatorId)
	if err != nil {
		panic(err)
	}
	if !validOperator(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证被添加用户是否存在
	user, err := queryUserByUname(req.UserName)
	if err != nil {
		panic(err)
	}
	if existedUser(user) {
		WriteData(w, config.UserHasAlreadyExists)
		return
	}

	err = addUserData(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		content := "添加用户[" + req.UserName + "]"
		InsertOperateLog(operator.UserId.Hex(), operator.AgencyId, content, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 &&
		strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	userId := r.URL.Query().Get("user_id")

	if !util.ValidDeleteUser(operatorId, userId) {
		WriteData(w, config.InvalidParameterValue)
		return
	}

	// 验证操作人信息是否有效
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !validOperator(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证被删除的用户是否存在
	user, err := queryUserByID(userId)
	if err != nil {
		panic(err)
	}
	if !existedUser(user) {
		WriteData(w, config.UserHasNotExists)
		return
	}

	status := verifyOperatorPermission(operator, user.AgencyId)
	if status != config.Success {
		WriteData(w, config.NewError(status))
		return
	}
	
	
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func FetchUserList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func UpdatePwd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}


////=========== Private Methods ===========

// 根据用户ID获取用户信息
func queryUserByID(userId string) (model.User, error) {
	var user model.User
	objId := bson.ObjectIdHex(userId)
	query := func(c *mgo.Collection) error {
		return c.FindId(objId).One(&user)
	}
	err := SharedQuery(T_USER, query)
	return user, err
}

// 添加用户数据
func addUserData(req model.UserReq) error {
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"username":   req.UserName,
			"password":   util.MD5Encrypt(req.Password),
			"gender":     req.Gender,
			"birth":      req.Birth,
			"mobile":     req.Mobile,
			"agentid":    req.AgencyId,
			"role":       req.Role,
			"priority":   req.Priority,
			"lasttime":   0,
			"lastonip":   "",
			"status":     config.USER_STATUS_NORMAL,
			"createtime": time.Now().Unix(),
			"updatetime": time.Now().Unix(),
		}
		return c.Insert(selector)
	}
	return SharedQuery(T_USER, query)
}

// 删除用户信息
func deleteUserData()  {
	
}

// 验证操作员信息是否有效
func validOperator(operator model.User) bool {
	return operator.Status == config.USER_STATUS_NORMAL && len(operator.UserId) > 0 && len(operator.UserName) > 0
}

// 验证是否存在该用户
func existedUser(user model.User) bool {
	return user.Status > config.USER_STATUS_INVALID && len(user.UserId) > 0 && len(user.UserName) > 0
}

// 验证操作人的权限
func verifyOperatorPermission(operator model.User, agencyId string) int64 {
	if operator.Role == "customer" {
		return config.PermissionDeniedUser
	} else if operator.Role == "admin" {
		if len(agencyId) > 0 && agencyId == operator.AgencyId {
			return config.Success
		} else {
			return config.PermissionDeniedUser
		}
	}
	return config.Success
}
