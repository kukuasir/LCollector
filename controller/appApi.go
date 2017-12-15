package controller

import (
	"LCollector/config"
	"LCollector/model"
	"net/http"
	"strings"
)

func GetAppInfo(w http.ResponseWriter, r *http.Request) {

	if strings.Compare(r.Method, "OPTIONS") == 0 {
		WriteData(w, config.Success)
		return
	}

	var appRet model.AppRet
	appRet.ResultInfo.Status = config.Success
	appRet.ResultInfo.Message = config.TIPS_QUERY_SUCCEED
	appRet.AppData = config.App
	WriteData(w, appRet)
}
