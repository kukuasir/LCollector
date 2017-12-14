package model

import "gopkg.in/mgo.v2/bson"

/** 定义操作类型 */
const (
	OPERATE_TYPE_ADD    = 1
	OPERATE_TYPE_DELETE = 2
	OPERATE_TYPE_UPDATE = 3
)

/** 操作对象 */
const (
	OPERATE_TARGET_USER     = 1
	OPERATE_TARGET_AGENCY   = 2
	OPERATE_TARGET_DEVICE   = 3
	OPERATE_TARGET_PASSWORD = 4
)

/** 消息类型 */
const (

	/**
	1- 设备正常运行过程中，根据config中的"bh(s)"的时长来心跳，服务器在2倍的时间内，没有接收心跳信息，被认为掉线
	2- 服务器回复心跳可以是任何信息，如request,request,config等等，也可以没有任何消息
	*/
	MESSAGE_TYPE_HEARTBEAT = 1

	/**
	1- 设备正常运行过程中，根据config中的"rT1"和"rT2"设置的时间点来确定每天上报运行状态和配置信息
	2- 设备上线时，会主动上报一次状态信息，这个状态信息，就是作为config的基础
	*/
	MESSAGE_TYPE_STATUS = 2

	/**
	1- 基于设备最新的status字段信息对设备进行重新配置(服务端发送到设备)
	2- 配置信息发出后，设备会回复status，以确保配置是否正确
	*/
	MESSAGE_TYPE_CONFIG = 3

	/**
	1- 设备正常运行过程中，根据config中uItvl(m)设置的时间间隔来上报data数据，可能包含多组数据
	2- 当设备接收到一个request时，request中的"idx"字段包含了current，设备立即回复一次当前的数据，这个数据需要保存到数据库中
	*/
	MESSAGE_TYPE_DATA = 4

	/**
	1- 设备正常运行过程中，如果有超限数据，会立即上报一个warning数据
	2- 后面，根据config中wUItvl(m)设置的时间间隔来上报warning数据，直到采集数据正常
	3- 服务器接收到warning数据以后，需要短信提醒用户，短息提醒3次不再提醒，每个设备最多可以添加5个手机号码作为报警目标
	*/
	MESSAGE_TYPE_WARNING = 5

	/**
	这个是一个及时通讯接口，对设备进行远程控制
	*/
	MESSAGE_TYPE_REQUEST = 6
)

type LoginLog struct {
	UserId     bson.ObjectId `json:"user_id" bson:"user_id"`         // 用户ID
	UserName   string        `json:"user_name" bson:"user_name"`     // 用户名
	StatusDesc string        `json:"status_desc"`                    // 用户状态
	AgencyId   bson.ObjectId `json:"agency_id" bson:"agency_id"`     // 用户所属机构ID
	CreateTime int64         `json:"create_time" bson:"create_time"` // 创建时间
	SourceIP   string        `json:"source_ip" bson:"source_ip"`     // 来源IP
}

type TempLoginLog struct {
	UserId     bson.ObjectId `json:"user_id" bson:"user_id"`         // 用户ID
	Status     int64         `json:"status"`                         // 用户状态
	AgencyId   bson.ObjectId `json:"agency_id" bson:"agency_id"`     // 所属组织机构ID
	CreateTime int64         `json:"create_time" bson:"create_time"` // 创建时间
	SourceIP   string        `json:"source_ip" bson:"source_ip"`     // 来源IP
	UserNames  []string      `json:"user_names" bson:"user_names"`   // 用户名列表
}

type OperateLog struct {
	Type         int64         `json:"type"`                           // 操作类型(1:添加 2:删除 3:修改)
	Target       int64         `json:"target"`                         // 操作对象(1:用户 2:组织机构 3:设备 4:密码)
	OperatorId   bson.ObjectId `json:"operator_id" bson:"operator_id"` // 操作人ID
	OperatorName string        `json:"operator_name"`                  // 操作人名字
	AgencyId     bson.ObjectId `json:"agency_id" bson:"agency_id"`     // 用户所属机构ID
	TargetObject string        `json:"object", bson:"object"`          // 操作对象的ID或名称
	CreateTime   int64         `json:"create_time" bson:"create_time"` // 创建时间
	SourceIP     string        `json:"source_ip" bson:"source_ip"`     // 来源IP
}

type TempOperateLog struct {
	Type         int64         `json:"type"`                           // 操作类型
	Target       int64         `json:"target"`                         // 操作对象
	TargetObject string        `json:"object", bson:"object"`          // 操作对象的ID或名称
	OperatorId   bson.ObjectId `json:"operator_id" bson:"operator_id"` // 操作人ID
	AgencyId     bson.ObjectId `json:"agency_id" bson:"agency_id"`     // 用户所属机构ID
	CreateTime   int64         `json:"create_time" bson:"create_time"` // 创建时间
	SourceIP     string        `json:"source_ip" bson:"source_ip"`     // 来源IP
	UserNames    []string      `json:"user_names" bson:"user_names"`   // 用户名列表
}

type MessageLog struct {
	Type       int64         `json:"type"`                           // 消息类型
	DeviceId   bson.ObjectId `json:"device_id" bson:"device_id"`     // 发送或接收消息的设备编号
	Content    string        `json:"content"`                        // 发送或接收的消息内容
	CreateTime int64         `json:"create_time" bson:"create_time"` // 创建时间
	SourceIP   string        `json:"source_ip" bson:"source_ip"`     // 来源IP
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
