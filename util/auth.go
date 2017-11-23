package util

func Auth(uname string, pwd string) bool {
	if len(uname) == 0 || len(pwd) == 0 {
		return false
	}
	return true
}
