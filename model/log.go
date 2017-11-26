package model

import "gopkg.in/mgo.v2/bson"

type OperateLog struct {
	OperatorId bson.ObjectId `json:"operator_id"` // 操作人ID
	AgencyId   bson.ObjectId `json:"agency_id"`   // 操作人所属组织机构
	AgencyName string        `json:"agency_name"` // 所属机构名称
	Content    string        `json:"content"`     // 操作内容
	Time       int64         `json:"time"`        // 操作时间
	OnIP       string        `json:"on_ip"`       // 操作的IP
}

type MessageLog struct {
	DeviceId   string        `json:"device_id"`   // 发送或接收消息的设备编号
	AgencyId   bson.ObjectId `json:"agency_id"`   // 操作人所属组织机构
	AgencyName string        `json:"agency_name"` // 所属机构名称
	Content    string        `json:"content"`     // 发送或接收的消息内容
	Time       int64         `json:"time"`        // 发送或接收消息的时间
	OnIP       string        `json:"on_ip"`       // 设备的IP
}

type OperateLogRet struct {
	ResultInfo  Result       `json:"result"` // 返回结果
	OperateList []OperateLog `json:"datas"`  // 操作日志列表信息
}

type MessageLogRet struct {
	ResultInfo  Result       `json:"result"` // 返回结果
	MessageList []MessageLog `json:"datas"`  // 消息日志列表信息
}
