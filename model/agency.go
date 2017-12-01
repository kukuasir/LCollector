package model

import "gopkg.in/mgo.v2/bson"

type Agency struct {
	AgencyId      bson.ObjectId `json:"agency_id" bson:"_id,omitempty"`       // 组织机构唯一ID
	AgencyName    string        `json:"agency_name" bson:"agency_name"`       // 机构名称
	ContactName   string        `json:"contact_name" bson:"contact_name"`     // 联系人
	ContactNumber string        `json:"contact_number" bson:"contact_number"` // 联系电话
	ContactAddr   string        `json:"contact_addr" bson:"contact_addr"`     // 联系地址
	Status        int64         `json:"status"`                               // 状态
	CreateTime    int64         `json:"create_time" bson:"create_time"`       // 创建时间
	UpdateTime    int64         `json:"update_time" bson:"update_time"`       // 最后更新时间
}

type AgencyReq struct {
	OperatorId    string `json:"operator_id"`    // 操作人员ID
	AgencyId      string `json:"agency_id"`      // 组织机构ID
	AgencyName    string `json:"agency_name"`    // 机构名称
	ContactName   string `json:"contact_name"`   // 联系人
	ContactNumber string `json:"contact_number"` // 联系电话
	ContactAddr   string `json:"contact_addr"`   // 联系地址
	Status        int64  `json:"status"`         // 状态
}

type AgencyRet struct {
	ResultInfo Result `json:"result"` // 返回结果
	AgencyInfo Agency `json:"data"`   // 组织机构信息
}

type AgencyListRet struct {
	ResultInfo Result   `json:"result"` // 返回结果
	AgencyList []Agency `json:"datas"`  // 组织机构列表
}
