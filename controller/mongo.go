package controller

import (
	"LCollector/config"
	"gopkg.in/mgo.v2"
)

/** 定义表名，与Mongo中表名对应 */
const (
	T_USER         = "t_user"         // 用户信息表
	T_DEVICE       = "t_device"       // 设备信息表
	T_AGENCY       = "t_agency"       // 组织机构信息表
	T_LOGIN_LOG    = "t_login_log"    // 登录日志表
	T_OPERATE_LOG  = "t_operate_log"  // 操作日志表
	T_MESSAGE_LOG  = "t_message_log"  // 消息日志表
	T_USER_TOKEN   = "t_user_token"   // 用户Token关联表
	T_ROLE_PATH    = "t_role_path"    // 角色权限关联表
)

var session *mgo.Session

/**
 * 公共方法，获取collection对象
 */
func SharedQuery(table string, s func(*mgo.Collection) error) error {
	sess := GlobalMgoSession()
	defer sess.Close()
	coll := sess.DB(config.Mongo.Database).C(table)
	return s(coll)
}

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func GlobalMgoSession() *mgo.Session {
	if session == nil {
		var err error
		session, err = mgo.Dial(config.Mongo.Host)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	session.SetMode(mgo.Monotonic, true)
	//最大连接池默认为4096
	return session.Clone()
}
