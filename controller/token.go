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

var TOKEN_PREFIX = "BTK"

func HandleToken(user model.User) string {
	token := GenerateToken(user.UserId.Hex())
	if len(token) > 0 {
		success := SaveToken(user, token)
		if success {
			return token
		}
		return ""
	}
	return ""
}

func GenerateToken(userId string) string {
	str := TOKEN_PREFIX + "_" + userId
	return util.MD5Encrypt(str)
}

/**
1、通过用户id查找出来的token进行对比是否一致,如果不一致则新增一条记录；
2、如果一致，再判断该token是否过期，如果过期再延长有效期
*/
func SaveToken(user model.User, token string) bool {

	/** 传入的token跟数据库中查询的结果一致的处理 */
	if len(user.Token) > 0 && strings.Compare(token, user.Token) == 0 {
		timestamp := time.Now().Unix()
		if (user.Expire + config.System.ValidSecs) < timestamp {
			updateTokenValid(user.UserId, timestamp+config.System.ValidSecs)
		}
	} else { /** 传入的token跟数据库中查询的结果不一致的处理 */
		updateToken(user.UserId, token)
	}

	return true
}

// 修改Token有效期
func updateTokenValid(userId bson.ObjectId, timestamp int64) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{"$set": bson.M{"expire": timestamp + config.System.ValidSecs}}
		return c.UpdateId(userId, update)
	}
	return SharedQuery(T_USER, query)
}

// 修改Token数据
func updateToken(userId bson.ObjectId, token string) error {
	query := func(c *mgo.Collection) error {
		update := bson.M{"$set": bson.M{"token": token, "expire": time.Now().Unix()}}
		return c.UpdateId(userId, update)
	}
	return SharedQuery(T_USER, query)
}
