/*festajuninaweb main web application program for festajuninaweb
// -----------------------------------------------
// .../src/festajuninaweb/festajuninaweb.go

02-Sep-2018
Agora utilisamos varios microservices
- FestaJuninaWeb - website running on port :1710
- Main   - API :1605
- Dishes - API :1610
- Orders - API :1620

Quero rodar from containers.

*/
package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	cachehandler "feijoadajusu/areas/cachehandler"
	"feijoadajusu/areas/commonstruct"
	"feijoadajusu/areas/ordershandler"
	"feijoadajusu/areas/security"
	"fmt"
	"html/template"
	"log"
	"net/http"
	// The Models are shared by WEB and API
	disheshandler "feijoadajusu/areas/disheshandler"
	helper "feijoadajusu/areas/helper"
	dishes "feijoadajusu/models"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

// Message our message object
type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

var mongodbvar commonstruct.DatabaseX

var clients []Client

// var credentials helper.Credentials

var db *sql.DB
var err error
var redisclient *redis.Client
var sysid string

// Looks after the main routing
//
func main() {

	// db, err = sql.Open("mysql", "daniel:oculos18@/festajunina")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err.Error())
	// }

	redisclient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// -------------------------------------------------
	sysid = helper.GetSYSID()

	fmt.Println("<><><><> sysid: " + sysid)
	// -------------------------------------------------

	loadreferencedatainredis()

	// Read variables from server
	//
	envirvar := helper.GetDBParmAllFromCache()

	fmt.Println(">>> Web Server: restauranteweb.exe running.")
	fmt.Println("Loading reference data in cache - Redis")

	mongodbvar.Location = envirvar.APIMongoDBLocation
	mongodbvar.Database = envirvar.APIMongoDBDatabase

	fmt.Println("FEIJOADA web server")
	fmt.Println("Running... Web Server Listening to :" + envirvar.WEBServerPort)
	fmt.Println("API DISHES Server: " + "http://localhost" + " Port: " + envirvar.MSAPIdishesPort)
	fmt.Println("API ORDERS Server: " + "http://localhost" + " Port: " + envirvar.MSAPIordersPort)
	fmt.Println("API SECURITY Server: " + "http://localhost" + " Port: " + envirvar.SecurityMicroservice)

	router := XNewRouter()

	// handle using the router mux
	//
	http.Handle("/", router) // setting router rule

	http.Handle("/html/", http.StripPrefix("/html", http.FileServer(http.Dir("./"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js"))))
	http.Handle("/ts/", http.StripPrefix("/ts", http.FileServer(http.Dir("./ts"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir("./fonts"))))
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./images"))))

	err := http.ListenAndServe(envirvar.WEBServerPort, nil) // setting listening port
	// err := http.ListenAndServe(envirvar.WEBServerPort, nil) // setting listening port
	if err != nil {
		//using the mux router
		log.Fatal("ListenAndServe: ", err)
	}

}

func loadreferencedatainredis() {
	variable := helper.Readfileintostruct()
	err = redisclient.Set(sysid+"APIMongoDBLocation", variable.APIMongoDBLocation, 0).Err()
	err = redisclient.Set(sysid+"APIMongoDBDatabase", variable.APIMongoDBDatabase, 0).Err()
	err = redisclient.Set(sysid+"Web.APIServer.Port", variable.APIAPIServerPort, 0).Err()
	err = redisclient.Set(sysid+"WEBServerPort", variable.WEBServerPort, 0).Err()
	err = redisclient.Set(sysid+"Web.APIServer.IPAddress", variable.APIAPIServerIPAddress, 0).Err()
	err = redisclient.Set(sysid+"Web.Debug", variable.WEBDebug, 0).Err()
	err = redisclient.Set(sysid+"RecordCurrencyTick", variable.RecordCurrencyTick, 0).Err()
	err = redisclient.Set(sysid+"RunningFromServer", variable.RunningFromServer, 0).Err()
	err = redisclient.Set(sysid+"AppFestaJuninaEnabled", variable.AppFestaJuninaEnabled, 0).Err()
	err = redisclient.Set(sysid+"AppBitcoinEnabled", variable.AppBitcoinEnabled, 0).Err()
	err = redisclient.Set(sysid+"AppBelnorthEnabled", variable.AppBelnorthEnabled, 0).Err()
	err = redisclient.Set(sysid+"Organisation", variable.Organisation, 0).Err()

	err = redisclient.Set(sysid+"MSAPImainIPAddress", variable.MSAPImainIPAddress, 0).Err()
	err = redisclient.Set(sysid+"MSAPIdishesIPAddress", variable.MSAPIdishesIPAddress, 0).Err()
	err = redisclient.Set(sysid+"MSAPIordersIPAddress", variable.MSAPIordersIPAddress, 0).Err()
	err = redisclient.Set(sysid+"SecurityMicroservice", variable.SecurityMicroservice, 0).Err()
	err = redisclient.Set(sysid+"SecurityMicroserviceURL", variable.SecurityMicroserviceURL, 0).Err()

	err = redisclient.Set(sysid+"MSAPImainPort", variable.MSAPImainPort, 0).Err()
	err = redisclient.Set(sysid+"MSAPIdishesPort", variable.MSAPIdishesPort, 0).Err()
	err = redisclient.Set(sysid+"MSAPIordersPort", variable.MSAPIordersPort, 0).Err()
}

func root(httpwriter http.ResponseWriter, req *http.Request) {

	_, credentials := security.ValidateTokenV2(redisclient, req)

	// error, credentials := security.ValidateTokenV2(redisclient, req)

	// if error == "NotOkToLogin" {
	// 	http.Redirect(httpwriter, req, "/login", 303)
	// 	return
	// }

	helper.HomePage(httpwriter, redisclient, credentials)

}

// ----------------------------------------------------------
// Security section
// ----------------------------------------------------------

func signupPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.SignupPage(httpresponsewriter, httprequest, redisclient, sysid)
}

func forgotPasswordPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.ForgotPasswordPage(httpresponsewriter, httprequest, redisclient, sysid)
}

func requestCode(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.RequestCode(httpresponsewriter, httprequest, redisclient, sysid)
}

func changePassword(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.ChangePassword(httpresponsewriter, httprequest, redisclient, sysid)
}

// userRolesShowPage is invoked from initial menu item
func userRolesShowPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.UserRolesShowPage(httpresponsewriter, httprequest, redisclient)
}

// userRolesShowPage is invoked from the button to get user details
func userRolesGetDetails(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.UserRolesGetDetails(httpresponsewriter, httprequest, redisclient, sysid)
}

func userRolesUpdate(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.UserRolesUpdate(httpresponsewriter, httprequest, redisclient, sysid)
}

func logoutPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.LogoutPage(httpresponsewriter, httprequest)
}

func loginPageV4(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.LoginPage(httpresponsewriter, httprequest, redisclient, sysid)
}

func instructions(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.Instructions(httpresponsewriter, httprequest, redisclient)

}

// ----------------------------------------------------------
// Orders section
// ----------------------------------------------------------

func saveordertosql(httpwriter http.ResponseWriter, req *http.Request) {
	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	ordershandler.SaveOrderToMySQL(httpwriter, redisclient, credentials, req, sysid)

}

// ----------------
// Anonymous
// ----------------

func orderlist(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	// _, credentials := security.ValidateTokenV2(redisclient, req)

	// If user is not ADMIN, show only users order

	ordershandler.ListV2(httpwriter, redisclient, credentials, sysid)
	// ordershandler.ListV3OnlyPlaced(httpwriter, redisclient, credentials)
}

func orderadddisplay(httpwriter http.ResponseWriter, req *http.Request) {
	// _, credentials := security.ValidateTokenV2(redisclient, req)

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	// The method below calls disheshandler.Listdishes(redisclient)
	// ... and it has to be like that.
	ordershandler.LoadDisplayForAdd(httpwriter, redisclient, credentials, sysid)
}

// orderadd
// is designed to place an order for a client that is logged on
//
func orderadd(httpwriter http.ResponseWriter, req *http.Request) {

	// _, _ = security.ValidateTokenV2(redisclient, req)

	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	ordershandler.Add(httpwriter, req, redisclient, sysid)
}

// orderclientadd
// is designed to place an order for an anonymous client
// it creates a dummy client
func orderclientadd(httpwriter http.ResponseWriter, req *http.Request) {

	// Find token
	// Get user ID
	_, credentials := security.ValidateTokenV2(redisclient, req)

	ordershandler.AddOrderClient(httpwriter, req, redisclient, credentials, sysid)
}

func orderviewdisplay(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	ordershandler.LoadDisplayForView(httpwriter, req, redisclient, credentials, sysid)
}
func ordercancel(httpwriter http.ResponseWriter, req *http.Request) {
	_, _ = security.ValidateTokenV2(redisclient, req)

	error2 := ordershandler.OrderisCancelled(httpwriter, req, redisclient, sysid)

	if error2 == "401 Order being served" {
		http.Redirect(httpwriter, req, "/errorpage", 301)
		return

	}

	orderid := req.URL.Query().Get("orderid")
	backto := "/orderviewdisplay?orderid=" + orderid
	http.Redirect(httpwriter, req, backto, 303)
	// http.Redirect(httpwriter, req, "/orderlist", 303)
}

// ----------------

func orderlistcompleted(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	// Only Admin
	//
	if credentials.IsAdmin == "Yes" {
		ordershandler.ListCompleted(httpwriter, redisclient, credentials, sysid)
	}

}

func orderliststatus(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	// Only Admin
	//
	if credentials.IsAdmin == "Yes" {
		ordershandler.ListStatus(req, httpwriter, redisclient, credentials, sysid)
	}

}

func ordersettoserving(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.StartServing(httpwriter, req, redisclient, sysid)

	orderid := req.URL.Query().Get("orderid")
	backto := "/orderviewdisplay?orderid=" + orderid
	http.Redirect(httpwriter, req, backto, 303)
	// http.Redirect(httpwriter, req, "/orderlist", 303)
}

func ordersettoready(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.OrderisReady(httpwriter, req, redisclient, sysid)

	orderid := req.URL.Query().Get("orderid")
	backto := "/orderviewdisplay?orderid=" + orderid
	http.Redirect(httpwriter, req, backto, 303)
	// http.Redirect(httpwriter, req, "/orderlist", 303)
}

func ordersettocompleted(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.OrderisCompleted(httpwriter, req, redisclient, sysid)

	// orderid := req.URL.Query().Get("orderid")
	// backto := "/orderviewdisplay?orderid=" + orderid
	// http.Redirect(httpwriter, req, backto, 303)

	http.Redirect(httpwriter, req, "/orderlist", 303)
}

func ordersettopaylater(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.OrderWillBePaidLater(httpwriter, req, redisclient, sysid)

	orderid := req.URL.Query().Get("orderid")
	backto := "/orderviewdisplay?orderid=" + orderid
	http.Redirect(httpwriter, req, backto, 303)
	// http.Redirect(httpwriter, req, "/orderlist", 303)
}

func orderStartServing(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.StartServing(httpwriter, req, redisclient, sysid)
}

// ----------------------------------------------------------
// Dishes section
// ----------------------------------------------------------

func dishlistpictures(httpwriter http.ResponseWriter, req *http.Request) {

	_, credentials := security.ValidateTokenV2(redisclient, req)

	disheshandler.ListPictures(httpwriter, redisclient, credentials)
}

func dishlist(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	disheshandler.List(httpwriter, redisclient, credentials, sysid)
}

func dishadddisplay(httpwriter http.ResponseWriter, req *http.Request) {

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	disheshandler.LoadDisplayForAdd(httpwriter)
}

func dishadd(httpwriter http.ResponseWriter, httprequest *http.Request) {

	// Retornar credentials e passar para a rotina Add below
	//
	error, _ := security.ValidateTokenV2(redisclient, httprequest)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, httprequest, "/login", 303)
		return
	}

	disheshandler.Add(httpwriter, httprequest, redisclient, sysid)
}

func dishupdatedisplay(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, httprequest)

	if error == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}

	disheshandler.LoadDisplayForUpdate(httpresponsewriter, httprequest, redisclient, credentials, sysid)
}

func dishupdate(httpwriter http.ResponseWriter, httprequest *http.Request) {

	// Retornar credentials e passar para a rotina Add below
	//
	error, _ := security.ValidateTokenV2(redisclient, httprequest)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, httprequest, "/login", 303)
		return
	}

	disheshandler.Update(httpwriter, httprequest, redisclient, sysid)

}

func dishdeletedisplay(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}
	disheshandler.LoadDisplayForDelete(httpresponsewriter, httprequest, redisclient, sysid)

}

func dishdelete(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	// ---------------------------------------------------------------------
	//        Security - Authorisation Check
	// ---------------------------------------------------------------------
	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}
	// ---------------------------------------------------------------------

	dishtodelete := dishes.Dish{}

	dishtodelete.Name = httprequest.FormValue("dishname") // This is the key, must be unique
	dishtodelete.Type = httprequest.FormValue("dishtype")
	dishtodelete.Price = httprequest.FormValue("dishprice")
	dishtodelete.GlutenFree = httprequest.FormValue("dishglutenfree")
	dishtodelete.DairyFree = httprequest.FormValue("dishdairyfree")
	dishtodelete.Vegetarian = httprequest.FormValue("dishvegetarian")
	dishtodelete.InitialAvailable = httprequest.FormValue("dishinitialavailable")
	dishtodelete.CurrentAvailable = httprequest.FormValue("dishcurrentavailable")

	ret := disheshandler.DishdeleteAPI(dishtodelete)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
		return
	}
}

func dishdeletemultiple(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	// ---------------------------------------------------------------------
	//        Security - Authorisation Check
	// ---------------------------------------------------------------------
	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}
	// ---------------------------------------------------------------------

	httprequest.ParseForm()

	// Get all selected records
	dishselected := httprequest.Form["dishes"]

	var numrecsel = len(dishselected)

	if numrecsel <= 0 {
		http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
		return
	}

	ret := commonstruct.Resultado{}

	ret = disheshandler.DishDeleteMultipleAPI(dishselected)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
		return
	}

	http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
	return

}

func showcache(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	// ---------------------------------------------------------------------
	//        Security - Authorisation Check
	// ---------------------------------------------------------------------
	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}
	// ---------------------------------------------------------------------

	// Cache from API
	cachehandler.List(httpresponsewriter, redisclient, sysid)

}

func errorpage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	// create new template
	var listtemplate = `
	{{define "listtemplate"}}
	{{ .Info.Name }}
	{{end}}
	`
	t, _ := template.ParseFiles("templates/error.html")
	t, _ = t.Parse(listtemplate)

	t.Execute(httpresponsewriter, listtemplate)
	return
}

// ---------------------------------------------------
// Websockets
// ---------------------------------------------------
// ---------------------------------------------------
// func rootws(httpwriter http.ResponseWriter, httprequest *http.Request) {
// 	flag.Parse()
// 	hub := newHub()
// 	go hub.run()
// 	serveWs(hub, httpwriter, httprequest)

// 	// Upgrade initial GET request to a websocket

// 	wsconn, err := upgrader.Upgrade(httpwriter, httprequest, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Register our new client
// 	clients[wsconn] = true

// 	fmt.Println("Client subscribed")

// }

func rootws(httpwriter http.ResponseWriter, httprequest *http.Request) {
	conn, err := websocket.Upgrade(httpwriter, httprequest, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(msg))
	}

}

func broadcastmsg(httpwriter http.ResponseWriter, httprequest *http.Request) {
	fmt.Println("broadcast message")
}

func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Message")

	for _, c := range clients {
		broadcastsingle(msg, c)
	}

	fmt.Fprintf(w, "Broadcasting %v", msg)
}

func addClientAndGreet(list []Client, client Client) []Client {
	clients = append(list, client)
	websocket.WriteJSON(client.conn, Message{"Server", "Welcome!"})
	return clients
}

func broadcastsingle(msg []byte, client Client) {
	fmt.Printf("Broadcasting %+v\n", msg)
	client.send <- msg
}
