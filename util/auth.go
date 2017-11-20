package util

func Auth(username string, password string) bool {
	if len(username) == 0 || len(password) == 0 {
		return false
	}
	return true
}
