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
		"POST",
		"/agency/add",
		AddAgency,
	},

	Route{
		"DeleteAgency",
		"GET",
		"/agency/delete",
		DeleteAgency,
	},

	Route{
		"EditAgency",
		"POST",
		"/agency/edit",
		EditAgency,
	},

	Route{
		"FetchAgencyList",
		"GET",
		"/agency/list",
		FetchAgencyList,
	},

	Route{
		"GetAgencyInfo",
		"GET",
		"/agency/view",
		GetAgencyInfo,
	},

	Route{
		"GetAppInfo",
		"GET",
		"/getAppInfo",
		GetAppInfo,
	},

	Route{
		"AddDevice",
		"POST",
		"/device/add",
		AddDevice,
	},

	Route{
		"DeleteDevice",
		"GET",
		"/device/delete",
		DeleteDevice,
	},

	Route{
		"EditDevice",
		"POST",
		"/device/edit",
		EditDevice,
	},

	Route{
		"FetchDeviceList",
		"GET",
		"/device/list",
		FetchDeviceList,
	},

	Route{
		"GetDeviceInfo",
		"GET",
		"/device/view",
		GetDeviceInfo,
	},

	Route{
		"FetchMessageLogList",
		"GET",
		"/log/message",
		FetchMessageLogList,
	},

	Route{
		"FetchOperateLogList",
		"GET",
		"/log/operate",
		FetchOperateLogList,
	},

	Route{
		"Login",
		"POST",
		"/login",
		Login,
	},

	Route{
		"AddUser",
		"POST",
		"/user/add",
		AddUser,
	},

	Route{
		"DeleteUser",
		"GET",
		"/user/delete",
		DeleteUser,
	},

	Route{
		"EditUser",
		"POST",
		"/user/edit",
		EditUser,
	},

	Route{
		"FetchUserList",
		"GET",
		"/user/list",
		FetchUserList,
	},

	Route{
		"GetUserInfo",
		"GET",
		"/user/view",
		GetUserInfo,
	},

	Route{
		"UpdatePwd",
		"POST",
		"/user/updatePwd",
		UpdatePwd,
	},
}
