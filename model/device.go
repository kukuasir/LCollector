package model

import "gopkg.in/mgo.v2/bson"

type Device struct {
	DeviceId   bson.ObjectId `json:"device_id" bson:"_id,omitempty"` // 设备ID
	DeviceName string        `json:"device_name" bson:"device_name"` // 设备名称
	AgencyId   string        `json:"agency_id" bson:"agency_id"`     // 所属机构ID
	Latitude   float64       `json:"latitude"`                       // 维度
	Longitude  float64       `json:"longitude"`                      // 经度
	Status     int64         `json:"status"`                         // 状态(-1:未分配  0:正常  1:报废)
	CreateTime int64         `json:"create_time" bson:"create_time"` // 创建时间
	UpdateTime int64         `json:"update_time" bson:"update_time"` // 最后更新时间
}

type DeviceReq struct {
	OperatorId string  `json:"operator_id"` // 操作人员的ID
	DeviceId   string  `json:"device_id"`   // 设备编号
	DeviceName string  `json:"device_name"` // 设备名称
	AgencyId   string  `json:"agency_id"`   // 所属机构ID
	Latitude   float64 `json:"latitude"`    // 维度
	Longitude  float64 `json:"longitude"`   // 经度
	Status     int64   `json:"status"`      // 状态
}

type DeviceRet struct {
	ResultInfo Result `json:"result"` // 返回结果
	DeviceData Device `json:"data"`   // 设备信息
}

type DeviceListRet struct {
	ResultInfo Result   `json:"result"` // 返回结果
	DeviceList []Device `json:"datas"`  // 设备列表
}
