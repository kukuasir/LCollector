package model

import "gopkg.in/mgo.v2/bson"

/** 设备信息 */
type Device struct {
	DeviceNo   string        `json:"device_no" bson:"device_no"`                // 设备编号
	DeviceName string        `json:"device_name" bson:"device_name"`            // 设备名称
	AgencyId   bson.ObjectId `json:"agency_id" bson:"agency_id"`                // 所属机构ID
	AgencyName string        `json:"agency_name,omitempty", bson:"agency_name"` // 所属组织机构名称
	Latitude   float64       `json:"latitude"`                                  // 维度
	Longitude  float64       `json:"longitude"`                                 // 经度
	Status     int64         `json:"status"`                                    // 状态(-1:未分配  0:正常  1:报废)
	StatusDesc string        `json:"status_desc,omitempty"`                     // 状态描述
	CreateTime int64         `json:"create_time" bson:"create_time"`            // 创建时间
	UpdateTime int64         `json:"update_time" bson:"update_time"`            // 最后更新时间
}

/** 临时-设备与组织机构关联信息 */
type TempDevice struct {
	DeviceNo    string        `json:"device_no" bson:"device_no"`       // 设备编号
	DeviceName  string        `json:"device_name" bson:"device_name"`   // 设备名称
	AgencyId    bson.ObjectId `json:"agency_id" bson:"agency_id"`       // 所属机构ID
	Latitude    float64       `json:"latitude"`                         // 维度
	Longitude   float64       `json:"longitude"`                        // 经度
	Status      int64         `json:"status"`                           // 状态(-1:未分配  0:正常  1:报废)
	CreateTime  int64         `json:"create_time" bson:"create_time"`   // 创建时间
	UpdateTime  int64         `json:"update_time" bson:"update_time"`   // 最后更新时间
	AgencyNames []string      `json:"agency_names" bson:"agency_names"` // 组织机构名列表
}

/** 设备 */
type DeviceCheck struct {
	DeviceNo   string `json:"device_no" bson:"device_no"`     // 设备编号
	DeviceName string `json:"device_name" bson:"device_name"` // 设备名称
	Check      bool   `json:"check"`                          // 是否选中
}

/** 用于添加或修改设备信息请求的结构体 */
type DeviceReq struct {
	OperatorId string  `json:"operator_id"` // 操作人员的ID
	DeviceNo   string  `json:"device_no"`   // 设备编号
	Token      string  `json:"token"`       // Token
	DeviceName string  `json:"device_name"` // 设备名称
	AgencyId   string  `json:"agency_id"`   // 所属机构ID
	Latitude   float64 `json:"latitude"`    // 维度
	Longitude  float64 `json:"longitude"`   // 经度
	Status     int64   `json:"status"`      // 状态
}

/** Api返回的设备列表 */
type DeviceListRet struct {
	ResultInfo Result   `json:"result"` // 返回结果
	DeviceList []Device `json:"datas"`  // 设备列表
}

/** Api返回的设备九宫格列表 */
type DeviceGridRet struct {
	ResultInfo Result     `json:"result"` // 返回结果
	DeviceGrid [][]Device `json:"datas"`  // 设备列表
}

/** Api返回的设备信息 */
type DeviceRet struct {
	ResultInfo Result `json:"result"` // 返回结果
	DeviceData Device `json:"data"`   // 设备信息
}
