package model

import "gopkg.in/mgo.v2/bson"

type LoginLog struct {
	UserId     bson.ObjectId `json:"user_id" bson:"user_id"`         // 用户ID
	Status     int64         `json:"status"`                         // 用户状态
	AgencyId   bson.ObjectId `json:"agency_id" bson:"agency_id"`     // 用户所属机构ID
	CreateTime int64         `json:"create_time" bson:"create_time"` // 创建时间
	SourceIP   string        `json:"srouce_ip" bson:"srouce_ip"`     // 来源IP
}

type OperateLog struct {
	OperateType   int64         `json:"type" bson:"type"`               // 操作类型
	OperateTarget int64         `json:"target" bson:"target"`           // 操作对象
	OperatorId    bson.ObjectId `json:"operator_id" bson:"operator_id"` // 操作人ID
	AgencyId      bson.ObjectId `json:"agency_id" bson:"agency_id"`     // 用户所属机构ID
	TargetObject  string        `json:"object", bson:"object"`          // 操作对象的ID或名称
	CreateTime    int64         `json:"create_time" bson:"create_time"` // 创建时间
	SourceIP      string        `json:"srouce_ip" bson:"srouce_ip"`     // 来源IP
}

type MessageLog struct {
	DeviceId   bson.ObjectId `json:"device_id" bson:"device_id"`     // 发送或接收消息的设备编号
	AgencyId   bson.ObjectId `json:"agency_id" bson:"agency_id"`     // 用户所属机构ID
	Content    string        `json:"content"`                        // 发送或接收的消息内容
	CreateTime int64         `json:"create_time" bson:"create_time"` // 创建时间
	SourceIP   string        `json:"srouce_ip" bson:"srouce_ip"`     // 来源IP
}

type LoginLogRet struct {
	ResultInfo Result     `json:"result"` // 返回结果
	LoginList  []LoginLog `json:"datas"`  // 登录日志列表
}

type OperateLogRet struct {
	ResultInfo  Result       `json:"result"` // 返回结果
	OperateList []OperateLog `json:"datas"`  // 操作日志列表
}

type MessageLogRet struct {
	ResultInfo  Result       `json:"result"` // 返回结果
	MessageList []MessageLog `json:"datas"`  // 消息日志列表
}
