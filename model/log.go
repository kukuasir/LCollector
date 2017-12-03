package model

type LoginLog struct {
	UserId   string `json:"user_id" bson:"user_id"`     // 用户ID
	Status   int64  `json:"status"`                     // 用户状态
	AgencyId string `json:"agency_id" bson:"agency_id"` // 用户所属机构ID
	Time     int64  `json:"time"`                       // 登录具体时间
	OnIP     string `json:"on_ip" bson:"on_ip"`         // 登录的IP
}

type OperateLog struct {
	OperateType int64  `json:"operate_type" bson:"operate_type"` // 操作类型
	OperatorId  string `json:"operator_id" bson:"operator_id"`   // 操作人ID
	AgencyId    string `json:"agency_id" bson:"agency_id"`       // 操作人所属组织机构
	Target      string `json:"target"`                           // 操作对象
	Time        int64  `json:"time"`                             // 操作时间
	OnIP        string `json:"on_ip" bson:"on_ip"`               // 操作的IP
}

type MessageLog struct {
	DeviceId   string `json:"device_id" bson:"device_id"`     // 发送或接收消息的设备编号
	AgencyId   string `json:"agency_id" bson:"agency_id"`     // 操作人所属组织机构
	AgencyName string `json:"agency_name" bson:"agency_name"` // 所属机构名称
	Content    string `json:"content"`                        // 发送或接收的消息内容
	Time       int64  `json:"time"`                           // 发送或接收消息的时间
	OnIP       string `json:"on_ip" bson:"on_ip"`             // 设备的IP
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
