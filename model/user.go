package model

import "gopkg.in/mgo.v2/bson"

type User struct {
	UserId        bson.ObjectId `json:"user_id" bson:"_id,omitempty"`           // 用户唯一ID
	UserName      string        `json:"user_name" bson:"user_name"`             // 用户名(不可重复)
	Password      string        `json:"-"`                                      // 用户密码
	Gender        int64         `json:"gender"`                                 // 性别(0:男 1:女)
	Birth         string        `json:"birth"`                                  // 出生年月
	Mobile        string        `json:"mobile"`                                 // 联系方式
	AgencyId      string        `json:"agency_id" bson:"agency_id"`             // 所属机构ID
	AgencyName    string        `json:"agency_name" bson:"agency_name"`         // 所属机构名称
	Role          string        `json:"role"`                                   // 角色
	Priority      string        `json:"priority"`                               // 设备查看方式的优先级
	OwnDevids     []string      `json:"own_devids" bson:"own_devids"`           // 所分配的设备列表
	Status        int64         `json:"status"`                                 // 状态
	LastLoginTime int64         `json:"last_login_time" bson:"last_login_time"` // 最后一次登录时间
	LastLoginIP   string        `json:"last_login_ip" bson:"last_login_ip"`     // 最后一次登录的IP
	CreateTime    int64         `json:"create_time" bson:"create_time"`         // 创建时间
	UpdateTime    int64         `json:"update_time" bson:"update_time"`         // 最后更新时间
}

type UserToken struct {
	UserId string `json:"user_id"` // 用户唯一ID
	Token  string `json:"token"`   // token
	Expire int64  `json:"expire"`  // 失效时间
}

type UserReq struct {
	OperatorId string   `json:"operator_id"` // 操作人员的ID
	UserId     string   `json:"user_id"`     // 目标用户ID
	UserName   string   `json:"user_name"`   // 用户账号
	Password   string   `json:"password"`    // 用户密码
	Gender     int64    `json:"gender"`      // 性别(0:男 1:女)
	Birth      string   `json:"birth"`       // 出生年月
	Mobile     string   `json:"mobile"`      // 联系方式
	AgencyId   string   `json:"agency_id"`   // 所属机构ID
	Role       string   `json:"role"`        // 角色
	Priority   string   `json:"priority"`    // 设备查看方式的优先级
	OwnDevids  []string `json:"own_devids"`  // 所分配的设备列表
	Status     int64    `json:"status"`      // 状态
}

type UserRet struct {
	ResultInfo Result `json:"result"` // 返回结果
	UserData   User   `json:"data"`   // 用户信息
}

type UserListRet struct {
	ResultInfo Result `json:"result"` // 返回结果
	UserList   []User `json:"datas"`  // 用户列表信息
}
