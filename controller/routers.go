package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"AddAgency",
		"OPTIONS",
		"/agency/add",
		AddAgency,
	},

	Route{
		"AddAgency",
		"POST",
		"/agency/add",
		AddAgency,
	},

	Route{
		"DeleteAgency",
		"OPTIONS",
		"/agency/delete",
		DeleteAgency,
	},

	Route{
		"DeleteAgency",
		"GET",
		"/agency/delete",
		DeleteAgency,
	},

	Route{
		"EditAgency",
		"OPTIONS",
		"/agency/edit",
		EditAgency,
	},

	Route{
		"EditAgency",
		"POST",
		"/agency/edit",
		EditAgency,
	},

	Route{
		"FetchAgencyList",
		"OPTIONS",
		"/agency/list",
		FetchAgencyList,
	},

	Route{
		"FetchAgencyList",
		"GET",
		"/agency/list",
		FetchAgencyList,
	},

	Route{
		"GetAgencyInfo",
		"OPTIONS",
		"/agency/view",
		GetAgencyInfo,
	},

	Route{
		"GetAgencyInfo",
		"GET",
		"/agency/view",
		GetAgencyInfo,
	},

	Route{
		"GetAppInfo",
		"OPTIONS",
		"/getAppInfo",
		GetAppInfo,
	},

	Route{
		"GetAppInfo",
		"GET",
		"/getAppInfo",
		GetAppInfo,
	},

	Route{
		"AddDevice",
		"OPTIONS",
		"/device/add",
		AddDevice,
	},

	Route{
		"AddDevice",
		"POST",
		"/device/add",
		AddDevice,
	},

	Route{
		"RegisterDevice",
		"POST",
		"/device/register",
		RegisterDevice,
	},

	Route{
		"DeleteDevice",
		"OPTIONS",
		"/device/delete",
		DeleteDevice,
	},

	Route{
		"DeleteDevice",
		"GET",
		"/device/delete",
		DeleteDevice,
	},

	Route{
		"EditDevice",
		"OPTIONS",
		"/device/edit",
		EditDevice,
	},

	Route{
		"EditDevice",
		"POST",
		"/device/edit",
		EditDevice,
	},

	Route{
		"FetchDeviceList",
		"OPTIONS",
		"/device/list",
		FetchDeviceList,
	},

	Route{
		"FetchDeviceList",
		"GET",
		"/device/list",
		FetchDeviceList,
	},

	Route{
		"GetDeviceInfo",
		"OPTIONS",
		"/device/view",
		GetDeviceInfo,
	},

	Route{
		"GetDeviceInfo",
		"GET",
		"/device/view",
		GetDeviceInfo,
	},

	Route{
		"FetchMessageLogList",
		"OPTIONS",
		"/log/message",
		FetchMessageLogList,
	},

	Route{
		"FetchMessageLogList",
		"GET",
		"/log/message",
		FetchMessageLogList,
	},

	Route{
		"FetchOperateLogList",
		"OPTIONS",
		"/log/operate",
		FetchOperateLogList,
	},

	Route{
		"FetchOperateLogList",
		"GET",
		"/log/operate",
		FetchOperateLogList,
	},

	Route{
		"FetchLoginLogList",
		"OPTIONS",
		"/log/login",
		FetchLoginLogList,
	},

	Route{
		"FetchLoginLogList",
		"GET",
		"/log/login",
		FetchLoginLogList,
	},

	Route{
		"Login",
		"OPTIONS",
		"/login",
		Login,
	},

	Route{
		"Login",
		"POST",
		"/login",
		Login,
	},

	Route{
		"AddUser",
		"OPTIONS",
		"/user/add",
		AddUser,
	},

	Route{
		"AddUser",
		"POST",
		"/user/add",
		AddUser,
	},

	Route{
		"DeleteUser",
		"OPTIONS",
		"/user/delete",
		DeleteUser,
	},

	Route{
		"DeleteUser",
		"GET",
		"/user/delete",
		DeleteUser,
	},

	Route{
		"EditUser",
		"OPTIONS",
		"/user/edit",
		EditUser,
	},

	Route{
		"EditUser",
		"POST",
		"/user/edit",
		EditUser,
	},

	Route{
		"FetchUserList",
		"OPTIONS",
		"/user/list",
		FetchUserList,
	},

	Route{
		"FetchUserList",
		"GET",
		"/user/list",
		FetchUserList,
	},

	Route{
		"GetUserInfo",
		"OPTIONS",
		"/user/view",
		GetUserInfo,
	},

	Route{
		"GetUserInfo",
		"GET",
		"/user/view",
		GetUserInfo,
	},

	Route{
		"UpdatePwd",
		"OPTIONS",
		"/user/updatePwd",
		UpdatePwd,
	},

	Route{
		"UpdatePwd",
		"POST",
		"/user/updatePwd",
		UpdatePwd,
	},
}
