package controller

import (
	"LCollector/config"
	"LCollector/model"
	"LCollector/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

var TOKEN_PREFIX string = "BTK"

func GenerateToken(uid string) string {
	str := TOKEN_PREFIX + "_" + uid
	return util.MD5Encrypt(str)
}

/**
1、通过用户id查找出来的token进行对比是否一致,如果不一致则新增一条记录；
2、如果一致，再判断该token是否过期，如果过期再延长有效期
*/
func SaveToken(uid string, token string) bool {

	var userToken model.UserToken
	err := queryTokenByUserID(uid, &userToken)
	if err != nil {
		return false
	}

	/** 传入的token跟数据库中查询的结果一致的处理 */
	if len(userToken.Token) > 0 && strings.Compare(token, userToken.Token) == 0 {
		timestamp := time.Now().Unix()
		if (userToken.Expire + config.System.ValidTimes) < timestamp {
			updateTokenValid(uid, timestamp+config.System.ValidTimes)
		}
	} else { /** 传入的token跟数据库中查询的结果不一致的处理 */
		insertToken(uid, token)
	}

	return true
}

// 根据用户ID获取对应的Token
func queryTokenByUserID(uid string, ut *model.UserToken) error {
	query := func(c *mgo.Collection) error {
		selector := bson.M{"user_id": uid}
		return c.Find(selector).One(&ut)
	}
	return SharedQuery(T_USER_TOKEN, query)
}

// 修改Token有效期
func updateTokenValid(uid string, timestamp int64) error {
	query := func(c *mgo.Collection) error {
		selector := bson.M{"user_id": uid}
		update := bson.M{"expire": timestamp + config.System.ValidTimes}
		return c.Update(selector, update)
	}
	return SharedQuery(T_USER_TOKEN, query)
}

// 插入一条Token数据
func insertToken(uid, token string) error {
	query := func(c *mgo.Collection) error {
		selector := bson.M{
			"user_id": uid,
			"token":   token,
			"expire":  time.Now().Unix(),
		}
		return c.Insert(selector)
	}
	return SharedQuery(T_USER_TOKEN, query)
}
