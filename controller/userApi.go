package controller

import (
	"LCollector/config"
	"LCollector/model"
	"LCollector/util"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddUser(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.UserReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.UserName) == 0 || len(req.Password) == 0 || len(req.Role) == 0 || len(req.AgencyId) == 0 {
		WriteData(w, config.InvalidParameterValue)
		return
	}

	// 验证操作人是否存在
	operator, err := queryUserByID(req.OperatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证被添加用户是否存在
	user, err := queryUserByUname(req.UserName)
	if err != nil {
		panic(err)
	}
	if ExistUser(user) {
		WriteData(w, config.UserHasAlreadyExists)
		return
	}

	err = addUserInfo(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		target := "用户[" + req.UserName + "]"
		InsertOperateLog(OPERATE_TYPE_ADD, operator, target, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	userId := r.URL.Query().Get("user_id")

	// 验证操作人是否存在
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证需要删除的用户是否存在
	user, err := queryUserByID(userId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(user) {
		WriteData(w, config.UserHasNotExists)
		return
	}

	// 验证操作人是否有权限删除对象
	status := verifyOperatorPermission(operator, user.Agency.AgencyId.Hex(), OPERATE_TARGET_USER)
	if status != config.Success {
		WriteData(w, config.NewError(status))
		return
	}

	err = deleteUserByID(user.UserId)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		target := "用户[" + user.UserId.Hex() + "]"
		InsertOperateLog(OPERATE_TYPE_DELETE, operator, target, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func EditUser(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.UserReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.UserName) == 0 || len(req.Password) == 0 || len(req.Role) == 0 || len(req.AgencyId) == 0 {
		WriteData(w, config.InvalidParameterValue)
		return
	}

	// 验证操作人是否存在
	operator, err := queryUserByID(req.OperatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证需要修改的用户是否存在
	user, err := queryUserByID(req.UserId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(user) {
		WriteData(w, config.UserHasNotExists)
		return
	}

	// 验证操作人是否有权限修改对象
	status := verifyOperatorPermission(operator, user.Agency.AgencyId.Hex(), OPERATE_TARGET_USER)
	if status != config.Success {
		WriteData(w, config.NewError(status))
		return
	}

	err = updateUserInfo(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		target := "用户[" + user.UserId.Hex() + "]"
		InsertOperateLog(OPERATE_TYPE_UPDATE, operator, target, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func UpdatePwd(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.UserReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.Password) == 0 {
		WriteData(w, config.InvalidParameterValue)
		return
	}

	// 验证操作人是否存在
	operator, err := queryUserByID(req.OperatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证用户是否存在
	user, err := queryUserByID(req.UserId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(user) {
		WriteData(w, config.UserHasNotExists)
		return
	}

	err = updatePwd(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		target := "密码[" + user.UserId.Hex() + "]"
		InsertOperateLog(OPERATE_TYPE_UPDATE, operator, target, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func FetchUserList(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if size == 0 {
		size = 20
	}

	// 验证操作人是否存在
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	userList, err := fetchPagingUserList(operator, page, size)
	if err != nil {
		panic(err)
	}

	// 返回查询结果
	var userListRet model.UserListRet
	userListRet.ResultInfo.Status = config.Success
	userListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	userListRet.UserList = userList
	WriteData(w, userListRet)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	userId := r.URL.Query().Get("user_id")

	// 验证操作人是否存在
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证需要查询的用户是否存在
	user, err := queryUserByID(userId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(user) {
		WriteData(w, config.UserHasNotExists)
		return
	}

	// 返回查询结果
	var userRet model.UserRet
	userRet.ResultInfo.Status = config.Success
	userRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	userRet.UserData = user
	WriteData(w, userRet)
}

////=========== Private Methods ===========

// 根据用户ID获取用户信息
func queryUserByID(userId string) (model.User, error) {
	var user model.User
	objId := bson.ObjectIdHex(userId)
	query := func(c *mgo.Collection) error {
		pipeline := []bson.M{
			bson.M{"$match": bson.M{"_id": objId}},
			bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
		}
		return c.Pipe(pipeline).One(&user)
	}
	err := SharedQuery(T_USER, query)
	return user, err
}

// 添加用户数据
func addUserInfo(req model.UserReq) error {
	query := func(c *mgo.Collection) error {
		insert := bson.M{
			"username":   req.UserName,
			"password":   util.MD5Encrypt(req.Password),
			"gender":     req.Gender,
			"birth":      req.Birth,
			"mobile":     req.Mobile,
			"agency_id":  bson.ObjectIdHex(req.AgencyId),
			"role":       req.Role,
			"priority":   req.Priority,
			"lasttime":   0,
			"lastonip":   "",
			"own_devids": req.OwnDevids,
			"status":     config.USER_STATUS_NORMAL,
			"createtime": time.Now().Unix(),
			"updatetime": time.Now().Unix(),
		}
		return c.Insert(insert)
	}
	return SharedQuery(T_USER, query)
}

// 删除用户信息(置为无效)
func deleteUserByID(userId bson.ObjectId) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{"$set": bson.M{"status": config.USER_STATUS_INVALID}}
		return c.UpdateId(userId, update)
	}
	return SharedQuery(T_USER, query)
}

// 修改用户信息
func updateUserInfo(req model.UserReq) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{
			"$set": bson.M{
				"gender":     req.Gender,
				"birth":      req.Birth,
				"mobile":     req.Mobile,
				"agency_id":  req.AgencyId,
				"role":       req.Role,
				"priority":   req.Priority,
				"own_devids": req.OwnDevids,
				"status":     req.Status,
				"updatetime": time.Now().Unix(),
			},
		}
		return c.UpdateId(req.UserId, update)
	}
	return SharedQuery(T_USER, query)
}

// 修改用户密码
func updatePwd(req model.UserReq) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{"$set": bson.M{"password": req.Password}}
		return c.UpdateId(req.UserId, update)
	}
	return SharedQuery(T_USER, query)
}

// 分页查询用户列表
func fetchPagingUserList(operator model.User, page, size int) ([]model.User, error) {
	var userList []model.User
	query := func(c *mgo.Collection) error {
		var pipeline []bson.M
		if operator.Role == "root" {
			pipeline = []bson.M{
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		} else {
			pipeline = []bson.M{
				bson.M{"$match": bson.M{"agency_id": operator.Agency.AgencyId}},
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		}
		return c.Pipe(pipeline).All(&userList)
	}
	err := SharedQuery(T_USER, query)
	return userList, err
}

// 验证操作人的权限
func verifyOperatorPermission(operator model.User, agencyId string, target int64) int64 {
	if operator.Role == "customer" {
		return errorWithTarget(target)
	} else if operator.Role == "admin" {
		if len(agencyId) > 0 && agencyId == operator.Agency.AgencyId.Hex() {
			return config.Success
		} else {
			return errorWithTarget(target)
		}
	}
	return config.Success
}

func errorWithTarget(target int64) int64 {
	switch target {
	case OPERATE_TARGET_USER:
		return config.PermissionDeniedUser
	case OPERATE_TARGET_AGENCY:
		return config.PermissionDeniedAgency
	case OPERATE_TARGET_DEVICE:
		return config.PermissionDeniedDevice
	default:
		return config.Success
	}
}

// 验证是否存在该用户
func ExistUser(user model.User) bool {
	return user.Status > config.USER_STATUS_INVALID && len(user.UserId) > 0
}
