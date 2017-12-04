package model

import "gopkg.in/mgo.v2/bson"

type LoginLog struct {
	UserId     bson.ObjectId `json:"user_id" bson:"user_id"`         // 用户ID
	Status     int64         `json:"status"`                         // 用户状态
	Agency     Agency        `json:"agency" bson:"agency"`           // 用户所属机构
	CreateTime int64         `json:"create_time" bson:"create_time"` // 创建时间
	SourceIP   string        `json:"srouce_ip" bson:"srouce_ip"`     // 来源IP
}

type OperateLog struct {
	OperateType int64         `json:"operate_type" bson:"operate_type"` // 操作类型
	OperatorId  bson.ObjectId `json:"operator_id" bson:"operator_id"`   // 操作人ID
	Agency      Agency        `json:"agency" bson:"agency"`             // 用户所属机构
	Target      string        `json:"target"`                           // 操作对象
	CreateTime  int64         `json:"create_time" bson:"create_time"`   // 创建时间
	SourceIP    string        `json:"srouce_ip" bson:"srouce_ip"`       // 来源IP
}

type MessageLog struct {
	DeviceId   bson.ObjectId `json:"device_id" bson:"device_id"`     // 发送或接收消息的设备编号
	Agency     Agency        `json:"agency" bson:"agency"`           // 用户所属机构
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
