package util

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
