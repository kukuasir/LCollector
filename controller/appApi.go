package controller

import (
	"LCollector/config"
	"LCollector/model"
	"net/http"
)

func GetAppInfo(w http.ResponseWriter, r *http.Request) {
	var appRet model.AppRet
	appRet.ResultInfo.Status = config.Success
	appRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	appRet.AppData = config.App
	WriteData(w, appRet)
}
