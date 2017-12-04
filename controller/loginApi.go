package controller

import (
	"LCollector/config"
	"LCollector/model"
	"LCollector/util"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "POST") != 0 && strings.Compare(r.Method, "OPTION") != 0 {
		WriteData(w, config.NewError(config.UnsupportedRequestMethod))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var loginReq model.LoginReq
	json.Unmarshal(body, &loginReq)

	username := loginReq.UserName
	password := util.MD5Encrypt(loginReq.Password)

	// 校验鉴权是否通过
	if !util.Auth(username, password) {
		WriteData(w, config.NewError(config.AuthenticateFailure))
		return
	}

	// 根据输入的用户名查询相应的用户信息
	user, err := queryUserByUname(username)
	if err != nil {
		panic(err)
	}
	if user.Status != config.USER_STATUS_NORMAL {
		WriteData(w, config.NewError(config.AccountHadBeenLocked))
		return
	}
	if strings.Compare(password, user.Password) != 0 {
		WriteData(w, config.NewError(config.InvalidAccountOrPassword))
		return
	}

	// 修改用户表中最后一次登录信息
	updateLastLoginInfo(user.UserId, r.RemoteAddr)

	// 记录到登录日志
	if config.Logger.EnableOperateLog {
		InsertLoginLog(user, r.RemoteAddr)
	}

	// 查询用户权限数据
	var paths []model.Path
	paths, err = queryUserResourcesByRole(user.Role)
	if err != nil {
		paths = defaultPaths()
	}

	loginret := newLoginRet(user, paths)
	WriteData(w, loginret)
}


////=========== Private Methods ===========

func queryUserByUname(uname string) (model.User, error) {
	var user model.User
	query := func(c *mgo.Collection) error {
		selector := bson.M{"user_name": uname}
		return c.Find(selector).One(&user)
	}
	err := SharedQuery(T_USER, query)
	return user, err
}

func updateLastLoginInfo(uid bson.ObjectId, ipAddr string) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{"$set": bson.M{"last_login_time": time.Now().Unix(), "last_login_ip": ipAddr}}
		return c.UpdateId(uid, update)
	}
	return SharedQuery(T_USER, query)
}

func queryUserResourcesByRole(role string) ([]model.Path, error) {
	var paths []model.Path
	query := func(c *mgo.Collection) error {
		selector := bson.M{"role": role}
		return c.Find(selector).All(&paths)
	}
	err := SharedQuery(T_ROLE_PATH, query)
	return paths, err
}

func defaultPaths() []model.Path {
	var paths []model.Path
	str := `[{"parent": "device", "children": ["name": "设备列表", "resource": "/deviceList"]},
			 {"parent": "system", "children": ["name": "操作日志", "resource": "/operateLog"]}]`
	json.Unmarshal([]byte(str), &paths)
	return paths
}

func newLoginRet(user model.User, paths []model.Path) model.LoginRet {
	var ret model.LoginRet
	ret.ResultInfo.Status = config.Success
	ret.ResultInfo.Message = config.TIPS_LOGIN_SUCCEED
	ret.ResultInfo.Token = generateToken(user.UserId.String())
	ret.RBACData.UserName = user.UserName
	ret.RBACData.Role = user.Role
	ret.RBACData.Paths = paths
	return ret
}

func generateToken(uid string) string {
	token := GenerateToken(uid)
	if len(token) > 0 {
		success := SaveToken(uid, token)
		if success {
			return token
		}
		return ""
	}
	return ""
}
