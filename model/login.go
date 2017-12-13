package model

type LoginReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginRet struct {
	ResultInfo Result `json:"result"` // 返回结果
	RBACData   RBAC   `json:"data"`   // 用户信息
}

type RBAC struct {
	OperatorId string `json:"operator_id"` // 用户id
	UserName   string `json:"user_name"`   // 用户名
	Role       string `json:"role"`        // 角色
	Paths      []Path `json:"paths"`       // 路径
}

type Path struct {
	Role     string     `json:"-"`
	Parent   string     `json:"parent"`   // 根节点
	Children []children `json:"children"` // 子节点
}

type children struct {
	Name     string `json:"name"`     // 节点中文名
	Resource string `json:"resource"` // 节点路径
}
