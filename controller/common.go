package controller

import (
	"LCollector/config"
	"LCollector/model"
	"encoding/json"
	"log"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"time"
)

func WriteData(w http.ResponseWriter, res interface{}) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Fatal("json marshal error: ", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Server", "BTK")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	defer func() {
		w.Write(data)
	}()
}

// 验证手机号格式
func CheckMobile(mobile string) bool {

	if len(mobile) != 11 {
		return false
	}

	/**
     * 手机号码:
     * 13[0-9], 14[5,7, 9], 15[0, 1, 2, 3, 5, 6, 7, 8, 9], 17[0-9], 18[0-9]
     */
	reg := `^1(3[0-9]|4[579]|5[0-35-9]|7[0-9]|8[0-9])\\d{8}$`
	rex := regexp.MustCompile(reg)
	return rex.MatchString(mobile)
}

// 验证Token是否有效
func ValidToken(user model.User, token string) bool {
	if !config.System.CheckToken {
		return true
	}
	if len(token) == 0 {
		return false
	}
	if token == user.Token {
		return false
	}
	if (user.Expire + config.System.ValidSecs) < time.Now().Unix() {
		return false
	} else {
		return true
	}
}

// 判断组织机构是否存在
func ExistAgency(agency model.Agency) bool {
	return len(agency.AgencyId) > 0 && agency.Status > config.AGENCY_STATUS_INVALID
}

// 判断设备是否存在
func ExistDevice(device model.Device) bool {
	return len(device.DeviceId) > 0 && device.Status > config.DEVICE_STATUS_INVALID
}

// 验证是否存在该用户
func ExistUser(user model.User) bool {
	return len(user.UserId) > 0 && user.Status > config.USER_STATUS_INVALID
}

// 查询列表的总个数
func GetCount(coll string) (int64, error) {
	var count int
	query := func(c *mgo.Collection) error {
		var err error
		count, err = c.Find(bson.M{}).Count()
		return err
	}
	err := SharedQuery(coll, query)
	return int64(count), err
}

// 校验page的值
func ValidPageValue(page *int) {
	if *page > 0 {
		*page = *page - 1
	} else {
		*page = 0
	}
}
