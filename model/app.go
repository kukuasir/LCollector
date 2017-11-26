package model

type AppInfo struct {
	Name      string `json:"name"`      // App名称
	Logo      string `json:"logo"`      // Logo
	Summary   string `json:"summary"`   // 描述
	Copyright string `json:"copyright"` // 版权
	QQ        string `json:"qq"`        // QQ号码
	Wechat    string `json:"wechat"`    // 微信公众号
	Website   string `json:"website"`   // 网站地址
	Version   string `json:"version"`   // 当前版本号
}
