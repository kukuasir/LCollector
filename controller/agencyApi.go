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

func AddAgency(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.AgencyReq
	json.Unmarshal(body, &req)

	// 验证数据是否有效
	//if len(req.AgencyName) == 0 || len(req.ContactName) == 0 || len(req.ContactNumber) == 0 || len(req.ContactAddr) == 0 {
	if len(req.AgencyName) == 0 {
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

	// 验证被添加的组织机构是否存在
	agency, err := queryAgencyInfoByName(req.AgencyName)
	if err != nil {
		panic(err)
	}
	if ExistAgency(agency) {
		WriteData(w, config.AgencyHasAlreadyExists)
		return
	}

	err = addAgencyInfo(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(OPERATE_TYPE_ADD, OPERATE_TARGET_AGENCY, operator, req.AgencyName, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func DeleteAgency(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	operatorId := r.URL.Query().Get("operator_id")
	agencyId := r.URL.Query().Get("agency_id")

	// 验证操作人是否存在
	operator, err := queryUserByID(operatorId)
	if err != nil {
		panic(err)
	}
	if !ExistUser(operator) {
		WriteData(w, config.OperaterHasNotExists)
		return
	}

	// 验证需要删除的机构是否存在
	agency, err := queryAgencyInfoByID(agencyId)
	if err != nil {
		panic(err)
	}
	if !ExistAgency(agency) {
		WriteData(w, config.AgencyHasNotExists)
		return
	}

	// 验证操作人是否有权限删除对象
	status := verifyOperatorPermission(operator, agencyId, OPERATE_TARGET_AGENCY)
	if status != config.Success {
		WriteData(w, config.NewError(status))
		return
	}

	err = deleteAgencyByID(agency.AgencyId)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(OPERATE_TYPE_DELETE, OPERATE_TARGET_DEVICE, operator, agencyId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func EditAgency(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req model.AgencyReq
	json.Unmarshal(body, &req)

	// 校验请求参数
	if len(req.AgencyName) == 0 {
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

	// 验证需要修改的机构是否存在
	agency, err := queryAgencyInfoByID(req.AgencyId)
	if err != nil {
		panic(err)
	}
	if !ExistAgency(agency) {
		WriteData(w, config.AgencyHasNotExists)
		return
	}

	// 验证操作人是否有权限修改对象
	status := verifyOperatorPermission(operator, req.AgencyId, OPERATE_TARGET_USER)
	if status != config.Success {
		WriteData(w, config.NewError(status))
		return
	}

	err = updateAgencyInfo(req)
	if err != nil {
		panic(err)
	}

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		InsertOperateLog(OPERATE_TYPE_UPDATE, OPERATE_TARGET_AGENCY, operator, req.AgencyId, r.RemoteAddr)
	}

	// 返回成功消息
	sucret := config.NewSuccess(config.TIPS_OPERA_SUCCEED)
	WriteData(w, sucret)
}

func FetchAgencyList(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	agencyList, err := fetchPagingAgencyList(page, size)
	if err != nil {
		panic(err)
	}

	// 返回查询结果
	var agencyListRet model.AgencyListRet
	agencyListRet.ResultInfo.Status = config.Success
	agencyListRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	agencyListRet.AgencyList = agencyList
	WriteData(w, agencyListRet)
}

func GetAgencyInfo(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "GET") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	agencyId := r.URL.Query().Get("agency_id")

	// 验证需要查询的用户是否存在
	agency, err := queryAgencyInfoByID(agencyId)
	if err != nil {
		panic(err)
	}
	if !ExistAgency(agency) {
		WriteData(w, config.AgencyHasNotExists)
		return
	}

	// 返回查询结果
	var agencyRet model.AgencyRet
	agencyRet.ResultInfo.Status = config.Success
	agencyRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	agencyRet.AgencyData = agency
	WriteData(w, agencyRet)
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
		return c.FindId(agencyId).One(&agency)
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
			"createtime":     time.Now().Unix(),
			"updatetime":     time.Now().Unix(),
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

// 修改组织机构信息
func updateAgencyInfo(req model.AgencyReq) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{
			"$set": bson.M{
				"agency_name":    req.AgencyName,
				"contact_name":   req.ContactName,
				"contact_number": req.ContactNumber,
				"contact_addr":   req.ContactAddr,
				"status":         req.Status,
				"updatetime":     time.Now().Unix(),
			},
		}
		return c.UpdateId(req.AgencyId, update)
	}
	return SharedQuery(T_AGENCY, query)
}

// 分页查询组织机构信息
func fetchPagingAgencyList(page, size int) ([]model.Agency, error) {
	var agencyList []model.Agency
	query := func(c *mgo.Collection) error {
		selector := bson.M{"status": bson.M{"$gt": config.AGENCY_STATUS_INVALID}}
		return c.Find(selector).Skip(page * size).Limit(size).All(&agencyList)
	}
	err := SharedQuery(T_AGENCY, query)
	return agencyList, err
}

// 判断组织机构是否存在
func ExistAgency(agency model.Agency) bool {
	if len(agency.AgencyId) == 0 {
		return false
	}
	return true
}
