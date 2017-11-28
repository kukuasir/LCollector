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
    "user_name": "######",
    "password": "######",
    "gender": NumberInt(######),
    "birth": "######",
    "mobile": "######",
    "agency_id": "######",
    "role": "root",
    "priority": "######",
    "status": NumberInt(0),
    "last_login_time": NumberLong(0),
    "last_login_ip": "",
    "create_time": NumberLong(######),
    "update_time": NumberLong(######)
});

db.t_agency.find({});
db.t_user.find({});
