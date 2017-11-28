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
	var user model.User
	err := queryUserByUserName(username, &user)
	if err != nil {
		panic(err)
	}
	if user.Status != config.USER_STATUS_NORMAL {
		WriteData(w, config.NewError(config.AccountHadBeenLocked))
	}
	if strings.Compare(password, user.Password) != 0 {
		WriteData(w, config.NewError(config.InvalidAccountOrPassword))
	}

	// 修改用户表中最后一次登录信息
	err = updateLastLoginInfo(user.UserId, r.RemoteAddr)

	// 记录操作日志
	if config.Logger.EnableOperateLog {
		info := "用户登录[" + user.UserId.String() + "]"
		InsertOperateLog(user.UserId.String(), user.AgencyId, info, r.RemoteAddr)
	}

	// 查询用户权限数据
	var paths []model.Path
	err = queryUserResourcesByRole(user.Role, paths)
	if err != nil {
		paths = defaultPaths()
	}

	loginret := newLoginRet(user, paths)
	WriteData(w, loginret)
}

func queryUserByUserName(uname string, user *model.User) error {
	query := func(c *mgo.Collection) error {
		selector := bson.M{"user_name": uname}
		return c.Find(selector).One(&user)
	}
	return SharedQuery(T_USER, query)
}

func updateLastLoginInfo(uid bson.ObjectId, ipAddr string) error {
	query := func(c *mgo.Collection) error {
		selector := bson.M{"user_id": uid}
		update := bson.M{"last_login_time": time.Now().Unix(), "last_login_ip": ipAddr}
		return c.Update(selector, update)
	}
	return SharedQuery(T_USER, query)
}

func queryUserResourcesByRole(role string, paths []model.Path) error {
	query := func(c *mgo.Collection) error {
		selector := bson.M{"role": role}
		return c.Find(selector).All(&paths)
	}
	return SharedQuery(T_ROLE_PATH, query)
}

func defaultPaths() []model.Path {
	var paths []model.Path
	str := `[{"parent": "device", "children": ["/ownerDevice"]}, {"parent": "system", "children": ["/operateLog"]}]`
	json.Unmarshal([]byte(str), &paths)
	return paths
}

func newLoginRet(user model.User, paths []model.Path) model.LoginRet {
	var ret model.LoginRet
	ret.ResultInfo.Status = config.Success
	ret.ResultInfo.Message = "登录成功"
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
