package model

type AppRet struct {
	ResultInfo Result  `json:"result"`
	AppData    interface{} `json:"data"`
}
