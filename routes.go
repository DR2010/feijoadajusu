// Routes are defined here
// -----------------------------------------------
// .../src/restauranteweb/routes.go
// -----------------------------------------------
package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route is
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes are
type Routes []Route

// VNewRouter is
func VNewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

// XNewRouter is
func XNewRouter() *mux.Router {

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

var routes = Routes{
	// ----------------------------------------------------------- Root
	Route{"Index", "GET", "/", root},
	// ----------------------------------------------------------- Error
	// Route{"errorpage", "GET", "/login", errorpage},
	// ----------------------------------------------------------- Security
	Route{"login", "GET", "/login", loginPageV4},
	Route{"instructions", "GET", "/instructions", instructions},
	Route{"login", "POST", "/login", loginPageV4},
	Route{"logout", "GET", "/logout", logoutPage},
	Route{"signup", "GET", "/signup", signupPage},
	Route{"signup", "POST", "/signup", signupPage},
	Route{"forgotpassword", "GET", "/forgotpassword", forgotPasswordPage},
	Route{"forgotpassword", "POST", "/forgotpassword", forgotPasswordPage},
	Route{"requestcode", "POST", "/requestcode", requestCode},
	Route{"changepassword", "POST", "/changepassword", changePassword},
	Route{"userrolesshowpage", "GET", "/userrolesshowpage", userRolesShowPage},
	Route{"userrolesgetdetails", "POST", "/userrolesgetdetails", userRolesGetDetails},
	Route{"userrolesgetdetails", "GET", "/userrolesgetdetails", userRolesGetDetails},
	Route{"userrolesupdate", "POST", "/userrolesupdate", userRolesUpdate},
	// ----------------------------------------------------------- Dishes
	Route{"dishlist", "GET", "/dishlist", dishlist},
	Route{"dishlistpictures", "GET", "/dishlistpictures", dishlistpictures},
	Route{"dishadddisplay", "POST", "/dishadddisplay", dishadddisplay},
	Route{"dishupdatedisplay", "POST", "/dishupdatedisplay", dishupdatedisplay},
	Route{"dishdeletedisplay", "POST", "/dishdeletedisplay", dishdeletedisplay},
	Route{"dishdeletemultiple", "POST", "/dishdeletemultiple", dishdeletemultiple},
	Route{"dishadd", "POST", "/dishadd", dishadd},
	Route{"dishupdate", "POST", "/dishupdate", dishupdate},
	Route{"dishdelete", "POST", "/dishdelete", dishdelete},
	// ----------------------------------------------------------- Activity
	Route{"activitylist", "GET", "/activitylist", activitylist},
	Route{"activityadddisplay", "POST", "/activityadddisplay", activityadddisplay},

	// -----------------------------------------------------------

	Route{"showcache", "GET", "/showcache", showcache},
	Route{"errorpage", "POST", "/errorpage", errorpage},
	// ----------------------------------------------------------- Order
	Route{"saveordertosqlX", "GET", "/saveordertosqlX", saveordertosql},
	Route{"orderlist", "GET", "/orderlist", orderlist},
	Route{"orderlist", "POST", "/orderlist", orderlist},
	Route{"orderlistcompleted", "GET", "/orderlistcompleted", orderlistcompleted},
	Route{"orderlistcompleted", "POST", "/orderlistcompleted", orderlistcompleted},
	Route{"orderliststatus", "GET", "/orderliststatus", orderliststatus},
	Route{"orderliststatus", "POST", "/orderliststatus", orderliststatus},
	Route{"orderadddisplay", "POST", "/orderadddisplay", orderadddisplay},
	Route{"orderadddisplay", "GET", "/orderadddisplay", orderadddisplay},
	Route{"orderadd", "POST", "/orderadd", orderadd},
	Route{"orderclientadd", "POST", "/orderclientadd", orderclientadd},
	Route{"ordersettoserving", "GET", "/ordersettoserving", ordersettoserving},
	Route{"ordersettoready", "GET", "/ordersettoready", ordersettoready},
	Route{"ordersettocompleted", "GET", "/ordersettocompleted", ordersettocompleted},
	Route{"ordersettopaylater", "GET", "/ordersettopaylater", ordersettopaylater},
	Route{"ordercancelX", "GET", "/ordercancelX", ordercancel},
	Route{"ordercancelX", "POST", "/ordercancelX", ordercancel},
	Route{"orderviewdisplay", "GET", "/orderviewdisplay", orderviewdisplay},
	// -----------------------------------------------------------
}
