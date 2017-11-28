package model

import "gopkg.in/mgo.v2/bson"

type User struct {
	UserId        bson.ObjectId `json:"user_id" bson:"_id, omitempty"`
	UserName      string        `json:"user_name"`
	Password      string        `json:"-"`
	Gender        int64         `json:"gender"`          // 性别(0:男 1:女)
	Birth         string        `json:"birth"`           // 出生年月
	Mobile        string        `json:"mobile"`          // 联系方式
	AgencyId      string        `json:"agency_id"`       // 所属机构ID
	Role          string        `json:"role"`            // 角色
	Priority      string        `json:"priority"`        // 设备查看方式的优先级
	Status        int64         `json:"status"`          // 状态
	LastLoginTime int64         `json:"last_login_time"` // 最后一次登录时间
	LastLoginIP   string        `json:"last_login_ip"`   // 最后一次登录的IP
	CreateTime    int64         `json:"create_time"`     // 创建时间
	UpdateTime    int64         `json:"update_time"`     // 最后更新时间
}

type UserToken struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`  // token
	Expire int64  `json:"expire"` // 失效时间
}

type UserReq struct {
	OperatorId bson.ObjectId `json:"operator_id"` // 操作人员的ID
	UserId     bson.ObjectId `json:"user_id"`     // 目标用户ID
	UserName   string        `json:"user_name"`   // 用户账号(唯一)
	Gender     int64         `json:"gender"`      // 性别(0:男 1:女)
	Birth      string        `json:"birth"`       // 出生年月
	Mobile     string        `json:"mobile"`      // 联系方式
	AgencyId   bson.ObjectId `json:"agency_id"`   // 所属机构ID
	Role       string        `json:"role"`        // 角色
	Priority   string        `json:"priority"`    // 设备查看方式的优先级
	Status     int64         `json:"status"`      // 状态
}

type UserRet struct {
	ResultData Result `json:"result"` // 返回结果
	UserData   User   `json:"data"`   // 用户信息
}

type UserListRet struct {
	ResultData Result `json:"result"` // 返回结果
	UserList   []User `json:"datas"`  // 用户列表信息
}
