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
	"time"
	"strings"
)

func AddUser(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
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
	operator, err := queryUserBaseInfo(req.OperatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, req.Token) {
		WriteData(w, config.InvalidToken)
		return
	}

	// 验证被添加用户是否存在
	user, err := queryUserByUname(req.UserName)
	if ExistUser(user) {
		WriteData(w, config.UserHasAlreadyExists)
		return
	}

	err = addUserInfo(req)
	if err != nil {
		panic(err)
	}

	if len(req.DeviceIds) > 0 {
		err = UpsertDeviceIsBy(req.UserId, req.DeviceIds)
		if err != nil {
			panic(err)
		}
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(model.OPERATE_TYPE_ADD, model.OPERATE_TARGET_USER, operator, req.UserName, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	userId := r.URL.Query().Get("user_id")
	token := r.URL.Query().Get("token")

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, token) {
		WriteData(w, config.InvalidToken)
		return
	}

	// 验证需要删除的用户是否存在
	user, err := queryUserBaseInfo(userId)
	if err != nil {
		panic(err)
	}

	// 验证操作人是否有权限删除对象
	status := verifyOperatorPermission(operator, user.AgencyId.Hex(), model.OPERATE_TARGET_USER)
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
		InsertOperateLog(model.OPERATE_TYPE_DELETE, model.OPERATE_TARGET_USER, operator, user.UserId.Hex(), r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func EditUser(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
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
	operator, err := queryUserBaseInfo(req.OperatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, req.Token) {
		WriteData(w, config.InvalidToken)
		return
	}

	// 验证需要修改的用户是否存在
	user, err := queryUserBaseInfo(req.UserId)
	if err != nil {
		panic(err)
	}

	// 验证操作人是否有权限修改对象
	status := verifyOperatorPermission(operator, user.AgencyId.Hex(), model.OPERATE_TARGET_USER)
	if status != config.Success {
		WriteData(w, config.NewError(status))
		return
	}

	err = updateUserInfo(req)
	if err != nil {
		panic(err)
	}

	if len(req.DeviceIds) > 0 {
		err = UpsertDeviceIsBy(req.UserId, req.DeviceIds)
		if err != nil {
			panic(err)
		}
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(model.OPERATE_TYPE_UPDATE, model.OPERATE_TARGET_USER, operator, user.UserId.Hex(), r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func UpdatePwd(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
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
	operator, err := queryUserBaseInfo(req.OperatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, req.Token) {
		WriteData(w, config.InvalidToken)
		return
	}

	// 验证用户是否存在
	user, err := queryUserBaseInfo(req.UserId)
	if err != nil {
		panic(err)
	}

	err = updatePwd(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(model.OPERATE_TYPE_UPDATE, model.OPERATE_TARGET_PASSWORD, operator, user.UserId.Hex(), r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func FetchUserList(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	token := r.URL.Query().Get("token")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size == 0 {
		size = 20
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, token) {
		WriteData(w, config.InvalidToken)
		return
	}

	tempUsers, err := fetchPagingUserList(operator, page, size)
	if err != nil {
		panic(err)
	}

	// 转换到User表中
	var userList []model.User
	for i := 0; i < len(tempUsers); i++ {
		temp := tempUsers[i]
		var user model.User
		user.UserId = temp.UserId
		user.UserName = temp.UserName
		user.Gender = temp.Gender
		user.AgencyId = temp.AgencyId
		user.Role = temp.Role
		user.Status = temp.Status
		user.LastLoginTime = temp.LastLoginTime
		user.LastLoginIP = temp.LastLoginIP
		user.CreateTime = temp.CreateTime
		user.UpdateTime = temp.UpdateTime
		if len(temp.AgencyNames) > 0 {
			user.AgencyName = temp.AgencyNames[0]
		}
		userList = append(userList, user)
	}

	var totalCount int64
	if page == 0 {
		totalCount, err = GetCount(T_USER)
	}

	// 返回查询结果
	var userListRet model.UserListRet
	userListRet.ResultInfo.Status = config.Success
	userListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	userListRet.ResultInfo.Total = totalCount
	userListRet.UserList = userList
	WriteData(w, userListRet)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	userId := r.URL.Query().Get("user_id")
	token := r.URL.Query().Get("token")

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, token) {
		WriteData(w, config.InvalidToken)
		return
	}

	// 验证需要查询的用户是否存在
	temp, err := fetchUserInfo(userId)
	if err != nil {
		panic(err)
	}

	var user model.User
	user.UserId = temp.UserId
	user.UserName = temp.UserName
	user.Gender = temp.Gender
	user.AgencyId = temp.AgencyId
	user.Role = temp.Role
	user.Status = temp.Status
	user.StatusDesc = config.UserStatusDesc(temp.Status)
	user.LastLoginTime = temp.LastLoginTime
	user.LastLoginIP = temp.LastLoginIP
	user.CreateTime = temp.CreateTime
	user.UpdateTime = temp.UpdateTime
	if len(temp.AgencyNames) > 0 {
		user.AgencyName = temp.AgencyNames[0]
	}

	if !ExistUser(user) {
		WriteData(w, config.UserHasNotExists)
		return
	}

	// 获取用户可操作的设备信息
	var deviceCheckList []model.DeviceCheck
	deviceCheckList, err = fetchDeviceCheckListBy(user)
	user.Devices = deviceCheckList

	// 返回查询结果
	var userRet model.UserRet
	userRet.ResultInfo.Status = config.Success
	userRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	userRet.UserData = user
	WriteData(w, userRet)
}

////=========== Private Methods ===========

// 根据用户ID获取用户信息
func queryUserBaseInfo(userId string) (model.User, error) {
	var user model.User
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"_id":    bson.ObjectIdHex(userId),
			"status": bson.M{"$gt": config.USER_STATUS_INVALID},
		}
		return c.Find(selector).One(&user)
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
			"status":     config.USER_STATUS_NORMAL,
			"createtime": time.Now().Unix(),
			"updatetime": time.Now().Unix(),
			"token":      "",
			"expire":     0,
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

	set := make(bson.M)
	set["status"] = req.Status
	set["updatetime"] = time.Now().Unix()
	if req.Gender != 0 {
		set["gender"] = req.Gender
	}
	if len(req.Birth) > 0 {
		set["birth"] = req.Birth
	}
	if len(req.Mobile) > 0 {
		set["mobile"] = req.Mobile
	}
	if len(req.AgencyId) > 0 {
		set["agency_id"] = req.AgencyId
	}
	if len(req.Role) > 0 {
		set["role"] = req.Role
	}
	if len(req.Priority) > 0 {
		set["priority"] = req.Priority
	}

	query := func(c *mgo.Collection) error {
		update := bson.M{
			"$set": set,
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

// 把用户可操作的设备ID列表插入到关联表中
func UpsertDeviceIsBy(userId string, deviceIds []string) error {
	query := func(c *mgo.Collection) error {
		selector := bson.M{"user_id": bson.ObjectIdHex(userId)}
		update := bson.M{"$set": bson.M{"device_ids": deviceIds}}
		_, err := c.Upsert(selector, update)
		return err
	}
	return SharedQuery(T_USER, query)
}

// 分页查询用户列表
func fetchPagingUserList(operator model.User, page, size int) ([]model.TempUser, error) {

	var tempUsers []model.TempUser

	if operator.Role == "customer" {
		return tempUsers, nil
	}

	ValidPageValue(&page)

	query := func(c *mgo.Collection) error {
		pipeline := []bson.M{
			bson.M{"$skip": page * size},
			bson.M{"$limit": size},
			bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency_docs"}},
			bson.M{"$project": bson.M{
				"_id":             1,
				"user_name":       1,
				"gender":          1,
				"agency_id":       1,
				"role":            1,
				"status":          1,
				"last_login_time": 1,
				"last_login_ip":   1,
				"create_time":     1,
				"update_time":     1,
				"agency_names":    "$agency_docs.agency_name",
			}},
		}
		if operator.Role == "admin" {
			pipeline = append(pipeline, bson.M{"$match": bson.M{"agency_id": operator.AgencyId}})
		}
		return c.Pipe(pipeline).All(&tempUsers)
	}
	err := SharedQuery(T_USER, query)
	return tempUsers, err
}

// 根据用户ID查询可操作的设备列表
func fetchDeviceCheckListBy(user model.User) ([]model.DeviceCheck, error) {

	// 1、先获取该用户机构下未被删除的所有设备
	totalDevices, err := fetchDeviceListInAgecy(user.AgencyId)

	// 2、再获取该用户可操作的设备列表
	deviceDocs, err := fetchDeviceListInUsed(user.UserId)

	// 3、对用户可使用的设备进行判断
	var deviceCheckList []model.DeviceCheck
	for i := 0; i < len(totalDevices); i++ {
		temp := totalDevices[i]
		var deviceCheck model.DeviceCheck
		deviceCheck.DeviceId = temp.DeviceId
		deviceCheck.DeviceName = temp.DeviceName
		deviceCheck.Check = false
		for j := 0; j < len(deviceDocs); j++ {
			usedDevice := deviceDocs[j].UsableDevices[0]
			if len(usedDevice.DeviceId) > 0 && usedDevice.Status > config.DEVICE_STATUS_INVALID {
				if usedDevice.DeviceId == temp.DeviceId {
					deviceCheck.Check = true
					break
				}
			}
		}
		deviceCheckList = append(deviceCheckList, deviceCheck)
	}

	return deviceCheckList, err
}

func fetchDeviceListInAgecy(agencyId bson.ObjectId) ([]model.Device, error) {
	var devices []model.Device
	query := func(c *mgo.Collection) error {
		selector := bson.M{"agency_id": agencyId, "status": bson.M{"$gt": config.DEVICE_STATUS_INVALID}}
		return c.Find(selector).All(&devices)
	}
	err := SharedQuery(T_DEVICE, query)
	return devices, err
}

func fetchDeviceListInUsed(userId bson.ObjectId) ([]model.TempUser, error) {
	var usedDevices []model.TempUser
	query := func(c *mgo.Collection) error {
		pipeline := []bson.M{
			bson.M{"$match": bson.M{"_id": userId}},
			bson.M{"$unwind": "$device_ids"},
			bson.M{"$lookup": bson.M{"from": T_DEVICE, "localField": "device_ids", "foreignField": "_id", "as": "device_docs"}},
			bson.M{"$project": bson.M{
				"_id":                     1,
				"device_docs._id":         1,
				"device_docs.device_name": 1,
				"device_docs.latitude":    1,
				"device_docs.longitude":   1,
				"device_docs.status":      1,
				"device_docs.create_time": 1,
				"device_docs.update_time": 1,
			}},
		}
		return c.Pipe(pipeline).All(&usedDevices)
	}
	err := SharedQuery(T_USER, query)
	return usedDevices, err
}

func fetchUserInfo(userId string) (model.TempUser, error) {
	var tempUser model.TempUser
	objId := bson.ObjectIdHex(userId)
	query := func(c *mgo.Collection) error {
		pipeline := []bson.M{
			bson.M{"$match": bson.M{"_id": objId}},
			bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency_docs"}},
			bson.M{"$project": bson.M{
				"_id":             1,
				"user_name":       1,
				"gender":          1,
				"agency_id":       1,
				"role":            1,
				"status":          1,
				"last_login_time": 1,
				"last_login_ip":   1,
				"create_time":     1,
				"update_time":     1,
				"agency_names":    "$agency_docs.agency_name",
			}},
		}
		return c.Pipe(pipeline).One(&tempUser)
	}
	err := SharedQuery(T_USER, query)
	return tempUser, err
}

// 验证操作人的权限
func verifyOperatorPermission(operator model.User, agencyId string, target int64) int64 {
	if operator.Role == "customer" {
		return errorWithTarget(target)
	} else if operator.Role == "admin" {
		if len(agencyId) > 0 && agencyId == operator.AgencyId.Hex() {
			return config.Success
		} else {
			return errorWithTarget(target)
		}
	}
	return config.Success
}

func errorWithTarget(target int64) int64 {
	switch target {
	case model.OPERATE_TARGET_USER:
		return config.PermissionDeniedUser
	case model.OPERATE_TARGET_AGENCY:
		return config.PermissionDeniedAgency
	case model.OPERATE_TARGET_DEVICE:
		return config.PermissionDeniedDevice
	default:
		return config.Success
	}
}
