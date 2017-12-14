# LCollector

Basic MVC Web Application in Go

## 第三方依赖
go get github.com/BurntSushi/toml  # 读取配置文件
go get github.com/mikespook/gorbac  # RBAC第三方库
go get gopkg.in/mgo.v2  # MongoDB驱动

# 添加管理员用户
use Collector;

db.t_agency.insert({
    "agency_name": "######",
    "contact_name": "######",
    "contact_number": "######",
    "contact_addr": "######",
    "status": NumberInt(0),
    "create_time": NumberLong(######),
    "update_time": NumberLong(######)
});

db.t_user.insert({
    #"_id" : ObjectId("5a1d14d246400555786d169a"),
    "user_name" : "root",
    "password" : "ff9830c42660c1dd1942844f8069b74a",
    "gender" : 1,
    "birth" : "1986-06-01",
    "mobile" : "13200000001",
    "agency_id" : ObjectId("5a1d11b946400555786d1699"),
    "role" : "root",
    "priority" : "123",
    "status" : 0,
    "last_login_time" : NumberLong(0),
    "last_login_ip" : "",
    "create_time" : NumberLong(1513266349),
    "update_time" : NumberLong(1513266349)
});

db.t_agency.find({});
db.t_user.find({});
