package model

import (
	"gopkg.in/mgo.v2/bson"
)

/** 用户信息 */
type User struct {
	UserId        bson.ObjectId `json:"user_id" bson:"_id,omitempty"`              // 用户唯一ID
	UserName      string        `json:"user_name" bson:"user_name"`                // 用户名(不可重复)
	Password      string        `json:"-"`                                         // 用户密码
	Gender        int64         `json:"gender"`                                    // 性别(0:男 1:女)
	Birth         string        `json:"birth,omitempty"`                           // 出生年月
	Mobile        string        `json:"mobile,omitempty"`                          // 联系方式
	AgencyId      bson.ObjectId `json:"agency_id" bson:"agency_id"`                // 所属组织机构ID
	AgencyName    string        `json:"agency_name,omitempty", bson:"agency_name"` // 所属组织机构名称
	Role          string        `json:"role"`                                      // 角色
	Priority      string        `json:"priority,omitempty"`                        // 设备查看方式的优先级
	Status        int64         `json:"status"`                                    // 状态
	StatusDesc    string        `json:"status_desc"`                               // 状态描述
	LastLoginTime int64         `json:"last_login_time" bson:"last_login_time"`    // 最后一次登录时间
	LastLoginIP   string        `json:"last_login_ip" bson:"last_login_ip"`        // 最后一次登录的IP
	CreateTime    int64         `json:"create_time" bson:"create_time"`            // 创建时间
	UpdateTime    int64         `json:"update_time" bson:"update_time"`            // 最后更新时间
	Devices       []DeviceCheck `json:"devices,omitempty"`                         // 可以查看或操作的设备列表
}

/** 临时-用户与组织机构关联信息 */
type TempUser struct {
	UserId        bson.ObjectId   `json:"user_id" bson:"_id,omitempty"`           // 用户唯一ID
	UserName      string          `json:"user_name" bson:"user_name"`             // 用户名(不可重复)
	Gender        int64           `json:"gender"`                                 // 性别(0:男 1:女)
	AgencyId      bson.ObjectId   `json:"agency_id" bson:"agency_id"`             // 所属组织机构ID
	Role          string          `json:"role"`                                   // 角色
	Status        int64           `json:"status"`                                 // 状态
	LastLoginTime int64           `json:"last_login_time" bson:"last_login_time"` // 最后一次登录时间
	LastLoginIP   string          `json:"last_login_ip" bson:"last_login_ip"`     // 最后一次登录的IP
	CreateTime    int64           `json:"create_time" bson:"create_time"`         // 创建时间
	UpdateTime    int64           `json:"update_time" bson:"update_time"`         // 最后更新时间
	AgencyNames   []string        `json:"agency_names" bson:"agency_names"`       // 组织机构名列表
	UsableDevices []Device        `json:"device_docs" bson:"device_docs"`         // 用户可操作的设备列表
}

/** 用于添加或修改用户信息请求的结构体 */
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
	Status     int64    `json:"status"`      // 状态
	DeviceIds  []string `json:"device_ids"`  // 可操作的设备ID列表
}

/** 用户可操作或查看的设备关联信息 */
type UserDevices struct {
	UserId  bson.ObjectId `json:"user_id" bson:"user_id"`
	Devices []Device      `json:"devices"`
}

/** 用户与Token关联信息 */
type UserToken struct {
	UserId bson.ObjectId `json:"user_id"` // 用户唯一ID
	Token  string        `json:"token"`   // token
	Expire int64         `json:"expire"`  // 失效时间
}

/** Api返回的用户列表 */
type UserListRet struct {
	ResultInfo Result `json:"result"` // 返回结果
	UserList   []User `json:"datas"`  // 用户列表信息
}

/** Api返回的用户信息 */
type UserRet struct {
	ResultInfo Result `json:"result"` // 返回结果
	UserData   User   `json:"data"`   // 用户信息
}
