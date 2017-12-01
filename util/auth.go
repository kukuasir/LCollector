package util

import (
	"LCollector/model"
)

const NAME_LENGTH_MIN int = 4
const NAME_LENGTH_MAX int = 12

/**
 验证账号和密码是否有效
 */
func Auth(uname string, pwd string) bool {

	nameLength := len(uname)
	pwdLength := len(pwd)

	if nameLength == 0 || pwdLength == 0 {
		return false
	} else if nameLength < NAME_LENGTH_MIN || nameLength > NAME_LENGTH_MAX {
		return false
	}
	return true
}

/**
 验证添加的用户信息是否有效
 */
func ValidAddUser(req model.UserReq) bool {
	if len(req.OperatorId) == 0 || len(req.UserName) == 0 || len(req.Password) == 0 || len(req.Role) == 0 || len(req.AgencyId) == 0 {
		return false
	}
	return true
}

/**
 验证需要删除的用户信息是否有效
 */
func ValidDeleteUser(operatorId string, userId string) bool {
	if len(operatorId) == 0 || len(userId) == 0 {
		return false
	}
	return true
}
