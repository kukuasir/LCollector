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
	DEVICE_STATUS_UNUSED  = 2  // 不可用
)

/** 组织机构状态 */
const (
	AGENCY_STATUS_INVALID = -1 // 无效的机构状态
	AGENCY_STATUS_NORMAL  = 0  // 正常
	AGENCY_STATUS_LOCKED  = 1  // 锁定
)
