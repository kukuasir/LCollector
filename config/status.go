package config

/** 用户状态 */
const (
	USER_STATUS_INVALID = -100 // 无效的用户状态
	USER_STATUS_LOCKED  = -1   // 锁定
	USER_STATUS_UNALLOC = 100  // 未分配
	USER_STATUS_NORMAL  = 200  // 正常
)

/** 设备状态 */
const (
	DEVICE_STATUS_INVALID = -100 // 无效的设备状态
	DEVICE_STATUS_LOCKED  = -1   // 锁定
	DEVICE_STATUS_UNALLOC = 100  // 未分配
	DEVICE_STATUS_NORMAL  = 200  // 正常
)

/** 组织机构状态 */
const (
	AGENCY_STATUS_INVALID = -100 // 无效的机构状态
	AGENCY_STATUS_LOCKED  = -1   // 锁定
	AGENCY_STATUS_NORMAL  = 200  // 正常
)

func UserStatusDesc(status int64) string {
	var desc string
	switch status {
	case USER_STATUS_INVALID:
		desc = "已删除"
		break
	case USER_STATUS_LOCKED:
		desc = "锁定"
		break
	case USER_STATUS_UNALLOC:
		desc = "未分配"
		break
	case USER_STATUS_NORMAL:
		desc = "正常"
		break
	default:
		desc = "正常"
		break
	}
	return desc
}

func AgencyStatusDesc(status int64) string {
	var desc string
	switch status {
	case AGENCY_STATUS_INVALID:
		desc = "已删除"
		break
	case AGENCY_STATUS_LOCKED:
		desc = "锁定"
		break
	case AGENCY_STATUS_NORMAL:
		desc = "正常"
		break
	default:
		desc = "正常"
		break
	}
	return desc
}

func DeviceStatusDesc(status int64) string {
	var desc string
	switch status {
	case DEVICE_STATUS_INVALID:
		desc = "已删除"
		break
	case DEVICE_STATUS_LOCKED:
		desc = "锁定"
		break
	case DEVICE_STATUS_UNALLOC:
		desc = "未分配"
		break
	case DEVICE_STATUS_NORMAL:
		desc = "正常"
		break
	default:
		desc = "正常"
		break
	}
	return desc
}
