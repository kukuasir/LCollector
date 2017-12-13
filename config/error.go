package config

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

/** 自定义错误码 */
const (
	UnsupportedRequestMethod = -100 // 不支持的方法

	AuthenticateFailure      = -1001 // 鉴权失败
	InvalidParameterValue    = -1002 // 数据格式有误
	AccountHadBeenLocked     = -1003 // 用户被锁定
	InvalidAccountOrPassword = -1004 // 无效的账号或密码

	DeviceHasAlreadyExists = -1010 // 设备已存在
	DeviceHasNotExists     = -1011 // 不存在的设备
	DeviceHasUnAlloc       = -1012 // 设备还未分配
	DeviceHasUnavailable   = -1013 // 设备不可用

	UserHasAlreadyExists = -1020 // 用户已经存在
	UserHasNotExists     = -1021 // 不存在的用户
	OperaterHasNotExists = -1022 // 不存在的操作人
	UserHadBeenLocked    = -1023 // 用户已被锁定

	AgencyHasAlreadyExists = -1030 // 组织机构已存在
	AgencyHasNotExists     = -1031 // 不存在的组织机构
	AgencyHadBeenLocked    = -1032 // 组织机构已被锁定

	PermissionDeniedDevice = -1040 // 没有设备的操作或查看权限
	PermissionDeniedUser   = -1041 // 没有用户的操作或查看权限
	PermissionDeniedAgency = -1042 // 没有机构的操作或查看权限

	Success = 200
)

/** 自定义成功提示 */
const (
	TIPS_LOGIN_SUCCEED = "登录成功"
	TIPS_QUERY_SUCCEED = "查询成功"
	TIPS_OPERA_SUCCEED = "操作成功"
)

type NSError struct {
	Code    int64
	Name    string
	Message string `json:"msg"`
}

var errors []NSError

func InitErrors() {

	data, err := ioutil.ReadFile("config/error.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &errors)
	if err != nil {
		panic(err)
	}
}

func NewSuccess(message string) *NSError {
	return &NSError{
		Code:    Success,
		Name:    "Success",
		Message: message,
	}
}

func NewError(code int64) *NSError {
	for _, error := range errors {
		if error.Code == code {
			return &error
		}
	}
	return NewError(Success)
}

func (e *NSError) Error() string {
	return "Code=" + strconv.FormatInt(e.Code, 10) + ", " + "Name=" + e.Name + ", " + "Message=" + e.Message
}
