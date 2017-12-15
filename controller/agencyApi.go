package controller

import (
	"LCollector/config"
	"LCollector/model"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
	"strings"
)

func AddAgency(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.AgencyReq
	json.Unmarshal(body, &req)

	// 验证数据是否有效
	//if len(req.AgencyName) == 0 || len(req.ContactName) == 0 || len(req.ContactNumber) == 0 || len(req.ContactAddr) == 0 {
	if len(req.AgencyName) == 0 {
		WriteData(w, config.NewError(config.InvalidParameterValue))
		return
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(req.OperatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, req.Token) {
		WriteData(w, config.NewError(config.InvalidToken))
		return
	}

	// 只有超级管理员才有权限操作组织机构
	if operator.Role != "root" {
		WriteData(w, config.NewError(config.PermissionDeniedAgency))
		return
	}

	// 验证被添加的组织机构是否存在
	agency, err := queryAgencyInfoByName(req.AgencyName)
	if ExistAgency(agency) {
		WriteData(w, config.NewError(config.AgencyHasAlreadyExists))
		return
	}

	err = addAgencyInfo(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(model.OPERATE_TYPE_ADD, model.OPERATE_TARGET_AGENCY, operator, req.AgencyName, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func DeleteAgency(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	agencyId := r.URL.Query().Get("agency_id")
	token := r.URL.Query().Get("token")

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, token) {
		WriteData(w, config.NewError(config.InvalidToken))
		return
	}

	// 只有超级管理员才有权限操作组织机构
	if operator.Role != "root" {
		WriteData(w, config.NewError(config.PermissionDeniedAgency))
		return
	}

	// 验证需要删除的机构是否存在
	agency, err := queryAgencyInfoByID(agencyId)
	if err != nil {
		panic(err)
	}

	// 删除该组织机构
	err = deleteAgencyByID(agency.AgencyId)
	if err != nil {
		panic(err)
	}

	// 同时把改组织机构下的用户都置为无效
	deleteUsersInAgency(agency.AgencyId, config.AGENCY_STATUS_INVALID)

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(model.OPERATE_TYPE_DELETE, model.OPERATE_TARGET_DEVICE, operator, agencyId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func EditAgency(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.AgencyReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.AgencyName) == 0 {
		WriteData(w, config.NewError(config.InvalidParameterValue))
		return
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(req.OperatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, req.Token) {
		WriteData(w, config.NewError(config.InvalidToken))
		return
	}

	// 只有超级管理员才有权限操作组织机构
	if operator.Role != "root" {
		WriteData(w, config.NewError(config.PermissionDeniedAgency))
		return
	}

	// 验证需要修改的机构是否存在
	_, err = queryAgencyInfoByID(req.AgencyId)
	if err != nil {
		panic(err)
	}

	err = updateAgencyInfo(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(model.OPERATE_TYPE_UPDATE, model.OPERATE_TARGET_AGENCY, operator, req.AgencyId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func FetchAgencyList(w http.ResponseWriter, r *http.Request) {

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
		WriteData(w, config.NewError(config.InvalidToken))
		return
	}

	agencyList, _ := fetchPagingAgencyList(page, size)

	var totalCount int64
	if page == 0 {
		totalCount, _ = GetCount(T_AGENCY)
	}

	// 返回查询结果
	var agencyListRet model.AgencyListRet
	agencyListRet.ResultInfo.Status = config.Success
	agencyListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	agencyListRet.ResultInfo.Total = totalCount
	agencyListRet.AgencyList = agencyList
	WriteData(w, agencyListRet)
}

func GetAgencyInfo(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	agencyId := r.URL.Query().Get("agency_id")
	token := r.URL.Query().Get("token")

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, token) {
		WriteData(w, config.NewError(config.InvalidToken))
		return
	}

	// 验证需要查询的用户是否存在
	agency, err := queryAgencyInfoByID(agencyId)
	if err != nil {
		panic(err)
	}

	agency.StatusDesc = config.AgencyStatusDesc(agency.Status)

	// 返回查询结果
	var agencyRet model.AgencyRet
	agencyRet.ResultInfo.Status = config.Success
	agencyRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	agencyRet.AgencyData = agency
	WriteData(w, agencyRet)
}

func FetchAgencyDevices(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	agencyId := r.URL.Query().Get("agency_id")
	token := r.URL.Query().Get("token")

	if len(agencyId) == 0 {
		panic("parameters error")
	}

	// 验证操作人是否存在
	operator, err := queryUserBaseInfo(operatorId)
	if err != nil {
		panic(err)
	}

	// 验证Token的有效性
	if !ValidToken(operator, token) {
		WriteData(w, config.NewError(config.InvalidToken))
		return
	}

	deviceList, _ := fetchDeviceListInAgency(agencyId)

	// 返回查询结果
	var deviceListRet model.DeviceListRet
	deviceListRet.ResultInfo.Status = config.Success
	deviceListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	deviceListRet.DeviceList = deviceList
	WriteData(w, deviceListRet)
}


////=========== Private Methods ===========

// 根据机构名称查询组织机构信息
func queryAgencyInfoByName(agencyName string) (model.Agency, error) {
	var agency model.Agency
	query := func(c *mgo.Collection) error {
		selector := bson.M{"agency_name": agencyName}
		return c.Find(selector).One(&agency)
	}
	err := SharedQuery(T_AGENCY, query)
	return agency, err
}

// 根据机构ID查询组织机构信息
func queryAgencyInfoByID(agencyId string) (model.Agency, error) {
	var agency model.Agency
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"_id": bson.ObjectIdHex(agencyId),
			"status": bson.M{"$gt": config.AGENCY_STATUS_INVALID},
		}
		return c.Find(selector).One(&agency)
	}
	err := SharedQuery(T_AGENCY, query)
	return agency, err
}

// 添加组织机构信息
func addAgencyInfo(req model.AgencyReq) error {
	query := func(c *mgo.Collection) error {
		insert := bson.M{
			"agency_name":    req.AgencyName,
			"contact_name":   req.ContactName,
			"contact_number": req.ContactNumber,
			"contact_addr":   req.ContactAddr,
			"status":         config.AGENCY_STATUS_NORMAL,
			"create_time":     time.Now().Unix(),
			"update_time":     time.Now().Unix(),
		}
		return c.Insert(insert)
	}
	return SharedQuery(T_AGENCY, query)
}

// 删除组织机构信息(状态设置为无效)
func deleteAgencyByID(agencyId bson.ObjectId) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{"$set": bson.M{"status": config.AGENCY_STATUS_INVALID}}
		return c.UpdateId(agencyId, update)
	}
	return SharedQuery(T_AGENCY, query)
}

// 修改该组织机构下的所有用户的状态
func deleteUsersInAgency(agencyId bson.ObjectId, status int64) error {
	query := func(c *mgo.Collection) error {
		selector := bson.M{"agency_id": agencyId}
		update := bson.M{"$set": bson.M{"status": status}}
		return c.Update(selector, update)
	}
	return SharedQuery(T_USER, query)
}

// 修改组织机构信息
func updateAgencyInfo(req model.AgencyReq) error {

	set := make(bson.M)
	set["status"] = req.Status
	set["update_time"] = time.Now().Unix()
	if len(req.AgencyName) > 0 {
		set["agency_name"] = req.AgencyName
	}
	if len(req.ContactName) > 0 {
		set["contact_name"] = req.ContactName
	}
	if len(req.ContactNumber) > 0 {
		set["contact_number"] = req.ContactNumber
	}
	if len(req.ContactAddr) > 0 {
		set["contact_addr"] = req.ContactAddr
	}

	query := func(c *mgo.Collection) error {
		update := bson.M{
			"$set": set,
		}
		return c.UpdateId(bson.ObjectIdHex(req.AgencyId), update)
	}
	return SharedQuery(T_AGENCY, query)
}

// 分页查询组织机构信息
func fetchPagingAgencyList(page, size int) ([]model.Agency, error) {

	// page值校验
	ValidPageValue(&page)

	// 分页查询
	var agencyList []model.Agency
	query := func(c *mgo.Collection) error {
		selector := bson.M{"status": bson.M{"$gt": config.AGENCY_STATUS_INVALID}}
		return c.Find(selector).Sort("-update_time").Skip(page * size).Limit(size).All(&agencyList)
	}
	err := SharedQuery(T_AGENCY, query)
	return agencyList, err
}

func fetchDeviceListInAgency(agencyId string) ([]model.Device, error) {
	var devlist []model.Device
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"agency_id": bson.ObjectIdHex(agencyId)}).Sort("-create_time").All(&devlist)
	}
	err := SharedQuery(T_DEVICE, query)
	return devlist, err
}
