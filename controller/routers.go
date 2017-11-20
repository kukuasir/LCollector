package controller

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
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
		"AgencyAddPost",
		"POST",
		"/agency/add",
		AgencyAddPost,
	},

	Route{
		"AgencyDeleteGet",
		"GET",
		"/agency/delete",
		AgencyDeleteGet,
	},

	Route{
		"AgencyEditPost",
		"POST",
		"/agency/edit",
		AgencyEditPost,
	},

	Route{
		"AgencyListGet",
		"GET",
		"/agency/list",
		AgencyListGet,
	},

	Route{
		"AgencyViewGet",
		"GET",
		"/agency/view",
		AgencyViewGet,
	},

	Route{
		"DeviceAddPost",
		"POST",
		"/device/add",
		DeviceAddPost,
	},

	Route{
		"DeviceDeleteGet",
		"GET",
		"/device/delete",
		DeviceDeleteGet,
	},

	Route{
		"DeviceEditPost",
		"POST",
		"/device/edit",
		DeviceEditPost,
	},

	Route{
		"DeviceListGet",
		"GET",
		"/device/list",
		DeviceListGet,
	},

	Route{
		"DeviceViewGet",
		"GET",
		"/device/view",
		DeviceViewGet,
	},

	Route{
		"LogMessageGet",
		"GET",
		"/log/message",
		LogMessageGet,
	},

	Route{
		"LogOperateGet",
		"GET",
		"/log/operate",
		LogOperateGet,
	},

	Route{
		"UserAddPost",
		"POST",
		"/user/add",
		UserAddPost,
	},

	Route{
		"UserDeleteGet",
		"GET",
		"/user/delete",
		UserDeleteGet,
	},

	Route{
		"UserEditPost",
		"POST",
		"/user/edit",
		UserEditPost,
	},

	Route{
		"UserListGet",
		"GET",
		"/user/list",
		UserListGet,
	},

	Route{
		"UserLoginPost",
		"POST",
		"/user/login",
		UserLoginPost,
	},

	Route{
		"UserUpdatePwdPost",
		"POST",
		"/user/updatePwd",
		UserUpdatePwdPost,
	},

	Route{
		"UserViewGet",
		"GET",
		"/user/view",
		UserViewGet,
	},

}