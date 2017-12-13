package config

/** 用户状态 */
const (
	USER_STATUS_INVALID = -1 // 无效的用户状态
	USER_STATUS_NORMAL  = 0  // 正常
	USER_STATUS_LOCKED  = 1  // 锁定
)

/** 设备状态 */
const (
	DEVICE_STATUS_INVALID = -1 // 无效的设备状态
	DEVICE_STATUS_UNALLOC = 0  // 未分配
	DEVICE_STATUS_NORMAL  = 1  // 正常
	DEVICE_STATUS_LOCKED  = 2  // 不可用
)

/** 组织机构状态 */
const (
	AGENCY_STATUS_INVALID = -1 // 无效的机构状态
	AGENCY_STATUS_NORMAL  = 0  // 正常
	AGENCY_STATUS_LOCKED  = 1  // 锁定
)

func UserStatusDesc(status int64) string {
	if status == USER_STATUS_NORMAL {
		return "正常"
	} else if status == USER_STATUS_LOCKED {
		return "锁定"
	} else {
		return "无效"
	}
}

func AgencyStatusDesc(status int64) string {
	if status == AGENCY_STATUS_NORMAL {
		return "正常"
	} else if status == AGENCY_STATUS_LOCKED {
		return "锁定"
	} else {
		return "无效"
	}
}

func DeviceStatusDesc(status int64) string {
	if status == DEVICE_STATUS_UNALLOC {
		return "未分配"
	} else if status == DEVICE_STATUS_NORMAL {
		return "正常"
	} else if status == DEVICE_STATUS_LOCKED {
		return "锁定"
	} else {
		return "无效"
	}
}
