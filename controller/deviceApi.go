package controller

import (
	"LCollector/config"
	"LCollector/model"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddDevice(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	var req model.DeviceReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.DeviceId) == 0 || len(req.DeviceName) == 0 {
		WriteData(w, config.InvalidParameterValue)
		return
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(req.OperatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 只有超级管理员才有权限添加设备
	if operator.Role != "root" {
		WriteData(w, config.NewError(config.PermissionDeniedDevice))
		return
	}

	// 验证需要添加的设备是否存在，存在则不能重复添加
	device, err := queryDeviceBaseInfo(req.DeviceId)
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
		InsertOperateLog(model.OPERATE_TYPE_ADD, model.OPERATE_TARGET_DEVICE, operator, req.DeviceId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func RegisterDevice(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 {
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
	device, err := queryDeviceBaseInfo(req.DeviceId)
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
		InsertMessageLog(model.MESSAGE_TYPE_STATUS, req.DeviceId, string(body), r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func DeleteDevice(w http.ResponseWriter, r *http.Request) {

	operatorId := r.URL.Query().Get("operator_id")
	deviceId := r.URL.Query().Get("device_id")

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 只有超级管理员才有权限删除设备
	if operator.Role != "root" {
		WriteData(w, config.NewError(config.PermissionDeniedDevice))
		return
	}

	// 验证需要删除的设备是否存在
	device, err := queryDeviceBaseInfo(deviceId)
	if err != nil {
		panic(err)
	}
	if !ExistDevice(device) {
		WriteData(w, config.DeviceHasNotExists)
		return
	}

	err = deleteDeviceByID(device.DeviceId)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(model.OPERATE_TYPE_DELETE, model.OPERATE_TARGET_DEVICE, operator, deviceId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func EditDevice(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	var req model.DeviceReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.DeviceName) == 0 || len(req.AgencyId) == 0 {
		WriteData(w, config.InvalidParameterValue)
		return
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(req.OperatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证需要修改的设备是否存在
	device, err := queryDeviceBaseInfo(req.DeviceId)
	if err != nil {
		panic(err)
	}
	if !ExistDevice(device) {
		WriteData(w, config.DeviceHasNotExists)
		return
	}

	// 验证操作人是否有权限修改对象
	status := verifyOperatorPermission(operator, device.AgencyId.Hex(), model.OPERATE_TARGET_DEVICE)
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
		InsertOperateLog(model.OPERATE_TYPE_UPDATE, model.OPERATE_TARGET_DEVICE, operator, req.DeviceId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func FetchDeviceList(w http.ResponseWriter, r *http.Request) {

	operatorId := r.URL.Query().Get("operator_id")
	if len(operatorId) == 0 {
		WriteData(w, config.NewError(config.InvalidParameterValue))
		return
	}

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
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	var deviceList []model.Device

	if operator.Role == "customer" {
		usedDevices, err := fetchDeviceListInUsed(operator.UserId)
		if err != nil {
			panic(err)
		}
		// 转换到Device表中
		for i := 0; i < len(usedDevices); i++ {
			temp := usedDevices[i].UsableDevices[0]
			var device model.Device
			device.DeviceId = temp.DeviceId
			device.DeviceName = temp.DeviceName
			device.Latitude = temp.Latitude
			device.Longitude = temp.Longitude
			device.CreateTime = temp.CreateTime
			device.UpdateTime = temp.UpdateTime
			deviceList = append(deviceList, device)
		}
	} else {
		tempDevices, err := fetchPagingDeviceList(operator, page, size)
		if err != nil {
			panic(err)
		}
		// 转换到Device表中
		for i := 0; i < len(tempDevices); i++ {
			temp := tempDevices[i]
			var device model.Device
			device.DeviceId = temp.DeviceId
			device.DeviceName = temp.DeviceName
			device.AgencyId = temp.AgencyId
			device.Latitude = temp.Latitude
			device.Longitude = temp.Longitude
			device.CreateTime = temp.CreateTime
			device.UpdateTime = temp.UpdateTime
			if len(temp.AgencyNames) > 0 {
				device.AgencyName = temp.AgencyNames[0]
			}
			deviceList = append(deviceList, device)
		}
	}

	var totalCount int64
	if page == 1 {
		totalCount, err = GetCount(T_AGENCY)
	}

	// 返回查询结果
	var deviceListRet model.DeviceListRet
	deviceListRet.ResultInfo.Status = config.Success
	deviceListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	deviceListRet.ResultInfo.Total = totalCount
	deviceListRet.DeviceList = deviceList
	WriteData(w, deviceListRet)
}

func GetDeviceInfo(w http.ResponseWriter, r *http.Request) {

	operatorId := r.URL.Query().Get("operator_id")
	deviceId := r.URL.Query().Get("device_id")

	if len(operatorId) == 0 || len(deviceId) == 0 {
		WriteData(w, config.NewError(config.InvalidParameterValue))
		return
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证需要查询的设备是否存在
	temp, err := fetchDeviceInfo(deviceId)
	if err != nil {
		panic(err)
	}

	var device model.Device
	device.DeviceId = temp.DeviceId
	device.DeviceName = temp.DeviceName
	device.AgencyId = temp.AgencyId
	device.Latitude = temp.Latitude
	device.Longitude = temp.Longitude
	device.CreateTime = temp.CreateTime
	device.UpdateTime = temp.UpdateTime
	if len(temp.AgencyNames) > 0 {
		device.AgencyName = temp.AgencyNames[0]
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

func queryDeviceBaseInfo(deviceId string) (model.Device, error) {
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

func fetchPagingDeviceList(operator model.User, page, size int) ([]model.TempDevice, error) {

	ValidPageValue(&page)

	var tempDevices []model.TempDevice
	query := func(c *mgo.Collection) error {
		pipeline := []bson.M{
			bson.M{"$skip": page * size},
			bson.M{"$limit": size},
			bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency_docs"}},
			bson.M{"$project": bson.M{
				"_id":          1,
				"device_name":  1,
				"agency_id":    1,
				"latitude":     1,
				"longitude":    1,
				"status":       1,
				"create_time":  1,
				"update_time":  1,
				"agency_names": "$agency_docs.agency_name",
			}},
		}
		if operator.Role == "admin" {
			pipeline = append(pipeline, bson.M{"$match": bson.M{"agency_id": operator.AgencyId}})
		}
		return c.Pipe(pipeline).All(&tempDevices)
	}
	err := SharedQuery(T_DEVICE, query)
	return tempDevices, err
}

func fetchDeviceInfo(deviceId string) (model.TempDevice, error) {
	var tempDevice model.TempDevice
	objId := bson.ObjectIdHex(deviceId)
	query := func(c *mgo.Collection) error {
		pipeline := []bson.M{
			bson.M{"$match": bson.M{"_id": objId}},
			bson.M{"$lookup": bson.M{"from": T_AGENCY, "localField": "agency_id", "foreignField": "_id", "as": "agency_docs"}},
			bson.M{"$project": bson.M{
				"_id":          1,
				"device_name":  1,
				"agency_id":    1,
				"latitude":     1,
				"longitude":    1,
				"status":       1,
				"create_time":  1,
				"update_time":  1,
				"agency_names": "$agency_docs.agency_name",
			}},
		}
		return c.Pipe(pipeline).One(&tempDevice)
	}
	err := SharedQuery(T_DEVICE, query)
	return tempDevice, err
}
