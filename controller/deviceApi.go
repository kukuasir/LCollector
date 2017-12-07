package controller

import (
	"LCollector/config"
	"LCollector/model"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"strconv"
)

func AddDevice(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.DeviceReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.DeviceId) == 0 || len(req.DeviceName) == 0 {
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

	// 验证需要添加的设备是否存在，存在则不能重复添加
	device, err := queryDeviceInfoByID(req.DeviceId)
	if err != nil {
		panic(err)
	}
	if ExistDevice(device) {
		WriteData(w, config.DeviceHasAlreadyExists)
		return
	}

	err = addDeviceInfo(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(OPERATE_TYPE_ADD, OPERATE_TARGET_DEVICE, operator, req.DeviceId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func RegisterDevice(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.DeviceReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.DeviceId) == 0 {
		WriteData(w, config.InvalidParameterValue)
		return
	}

	// 验证需要添加的设备是否存在，存在则不能重复添加
	device, err := queryDeviceInfoByID(req.DeviceId)
	if err != nil {
		panic(err)
	}

	if !ExistDevice(device) {
		err = addDeviceInfo(req)
	} else {
		err = updateDeviceInfo(req)
	}

	if err != nil {
		panic(err)
	}

	// 记录消息日志
	if config.Logger.EnableMessageLog {
		content := "注册设备[" + req.DeviceId + "]"
		InsertMessageLog(req.DeviceId, req.AgencyId, content, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func DeleteDevice(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	deviceId := r.URL.Query().Get("device_id")

	// 验证操作人是否存在
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证需要删除的设备是否存在
	device, err := queryDeviceInfoByID(deviceId)
	if err != nil {
		panic(err)
	}
	if !ExistDevice(device) {
		WriteData(w, config.DeviceHasNotExists)
		return
	}

	// 验证操作人是否有权限删除对象
	status := verifyOperatorPermission(operator, device.Agency.AgencyId.Hex(), OPERATE_TARGET_DEVICE)
	if status != config.Success {
		WriteData(w, config.NewError(status))
		return
	}

	err = deleteDeviceByID(device.DeviceId)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(OPERATE_TYPE_DELETE, OPERATE_TARGET_DEVICE, operator, deviceId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func EditDevice(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.DeviceReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.DeviceName) == 0 || len(req.AgencyId) == 0 {
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

	// 验证需要修改的设备是否存在
	device, err := queryDeviceInfoByID(req.DeviceId)
	if err != nil {
		panic(err)
	}
	if !ExistDevice(device) {
		WriteData(w, config.DeviceHasNotExists)
		return
	}

	// 验证操作人是否有权限修改对象
	status := verifyOperatorPermission(operator, device.Agency.AgencyId.Hex(), OPERATE_TARGET_DEVICE)
	if status != config.Success {
		WriteData(w, config.NewError(status))
		return
	}

	err = updateDeviceInfo(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(OPERATE_TYPE_UPDATE, OPERATE_TARGET_DEVICE, operator, req.DeviceId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func FetchDeviceList(w http.ResponseWriter, r *http.Request) {

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

	deviceList, err := fetchPagingDeviceList(operator, page, size)
	if err != nil {
		panic(err)
	}

	// 返回查询结果
	var deviceListRet model.DeviceListRet
	deviceListRet.ResultInfo.Status = config.Success
	deviceListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	deviceListRet.DeviceList = deviceList
	WriteData(w, deviceListRet)
}

func GetDeviceInfo(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	deviceId := r.URL.Query().Get("device_id")

	// 验证操作人是否存在
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证需要查询的设备是否存在
	device, err := queryDeviceInfoByID(deviceId)
	if err != nil {
		panic(err)
	}
	if !ExistDevice(device) {
		WriteData(w, config.DeviceHasNotExists)
		return
	}

	// 返回查询结果
	var deviceRet model.DeviceRet
	deviceRet.ResultInfo.Status = config.Success
	deviceRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	deviceRet.DeviceData = device
	WriteData(w, deviceRet)
}

////=========== Private Methods ===========

func queryDeviceInfoByID(deviceId string) (model.Device, error) {
	var device model.Device
	objId := bson.ObjectIdHex(deviceId)
	query := func(c *mgo.Collection) error {
		return c.FindId(objId).One(&device)
	}
	err := SharedQuery(T_DEVICE, query)
	return device, err
}

func addDeviceInfo(req model.DeviceReq) error {
	query := func(c *mgo.Collection) error {
		insert := bson.M{
			"_id":         req.DeviceId,
			"device_name": req.DeviceName,
			"agency_id":   req.AgencyId,
			"latitude":    req.Latitude,
			"longitude":   req.Longitude,
			"status":      config.DEVICE_STATUS_NORMAL,
			"createtime":  time.Now().Unix(),
			"updatetime":  time.Now().Unix(),
		}
		return c.Insert(insert)
	}
	return SharedQuery(T_DEVICE, query)
}

func deleteDeviceByID(deviceId bson.ObjectId) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{"$set": bson.M{"status": config.DEVICE_STATUS_INVALID}}
		return c.UpdateId(deviceId, update)
	}
	return SharedQuery(T_DEVICE, query)
}

func updateDeviceInfo(req model.DeviceReq) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{
			"$set": bson.M{
				"device_name": req.DeviceName,
				"agency_id":   req.AgencyId,
				"latitude":    req.Latitude,
				"longitude":   req.Longitude,
				"status":      req.Status,
				"updatetime":  time.Now().Unix(),
			},
		}
		return c.UpdateId(req.DeviceId, update)
	}
	return SharedQuery(T_DEVICE, query)
}

func fetchPagingDeviceList(operator model.User, page, size int) ([]model.Device, error) {
	var deviceList []model.Device
	query := func(c *mgo.Collection) error {
		var pipeline []bson.M
		if operator.Role == "root" {
			pipeline = []bson.M{
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"$unwind": "$agency"},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		} else if operator.Role == "admin" {
			pipeline = []bson.M{
				bson.M{"$match": bson.M{"agency_id": operator.Agency.AgencyId}},
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"$unwind": "$agency"},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		} else {
			pipeline = []bson.M{
				bson.M{"$match": bson.M{"agency_id": operator.Agency.AgencyId}},
				bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency"}},
				bson.M{"$unwind": "$agency"},
				bson.M{"$skip": page * size},
				bson.M{"$limit": size},
			}
		}
		return c.Pipe(pipeline).All(&deviceList)
	}
	err := SharedQuery(T_DEVICE, query)
	return deviceList, err
}

// 判断设备是否存在
func ExistDevice(device model.Device) bool {
	if len(device.DeviceId) == 0 {
		return false
	}
	return true
}
