package security

import (
	"encoding/json"
	"feijoadajusu/areas/commonstruct"
	"fmt"
	"io/ioutil"
	"log"
	"mongodb/dishes"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/template"
	"github.com/go-redis/redis"
)

// SignupPage is for the user to signup
func SignupPage(httpresponsewriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, sysid string) {

	type ControllerInfo struct {
		Name    string
		Message string
	}
	type DisplayTemplate struct {
		Info ControllerInfo
	}

	items := DisplayTemplate{}
	items.Info.Name = "Login Page"

	if req.Method != "POST" {

		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = ""
		t.Execute(httpresponsewriter, items)

		// http.ServeFile(res, req, "templates/security/signup.html")
		return
	}

	usernamemix := req.FormValue("username")
	preferredname := req.FormValue("preferredname")
	password := req.FormValue("password")
	passwordvalidate := req.FormValue("passwordvalidate")
	// applicationid := req.FormValue("applicationid")
	applicationid := "Restaurante" // Festa Junina

	username := strings.ToUpper(usernamemix)

	if username == "" {
		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = "Please enter Name."
		t.Execute(httpresponsewriter, items)
		return
	}

	if preferredname == "" {
		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = "Please enter Preferred Name."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == "" {
		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = "Please enter details."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == passwordvalidate {

		// Call API to check if user exists and create
		var resultado = SignUp(sysid, redisclient, username, preferredname, password, passwordvalidate, applicationid)
		if resultado.ErrorCode == "200 OK" {

		} else {
			// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
			t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
			items.Info.Message = "Passwords mismatch."
			t.Execute(httpresponsewriter, items)
			return
		}

		http.Redirect(httpresponsewriter, req, "/", 303)
	} else {
		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = "Passwords do not match."
		t.Execute(httpresponsewriter, items)
		return
	}

}

// UserRolesShowPage is for the user to signup
func UserRolesShowPage(httpresponsewriter http.ResponseWriter, req *http.Request, redisclient *redis.Client) {

	type ControllerInfo struct {
		Name          string
		Message       string
		UserName      string
		IsAdmin       string
		ApplicationID string
		Status        string
		UserType      string
	}
	type DisplayTemplate struct {
		Info ControllerInfo
	}

	items := DisplayTemplate{}
	items.Info.Name = "User Roles"

	// Call Server to get details
	//

	if req.Method != "POST" {

		t, _ := template.ParseFiles("templates/security/viewuserrolesheader.html", "templates/security/viewuserrolesdetails.html")
		items.Info.Message = ""
		t.Execute(httpresponsewriter, items)

		return
	}

	username := req.FormValue("username") // email address
	username = strings.ToUpper(username)

	if username == "" {
		t, _ := template.ParseFiles("templates/security/viewuserrolesheader.html", "templates/security/viewuserrolesdetails.html")
		items.Info.Message = "Please_enter_Name."
		t.Execute(httpresponsewriter, items)
		return
	}

}

// UserRolesGetDetails is for the user to signup
func UserRolesGetDetails(httpresponsewriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client, sysid string) {

	defer httprequest.Body.Close()
	bodybyte, _ := ioutil.ReadAll(httprequest.Body)

	ret := GetUserDetails(sysid, redisclient, bodybyte)

	obj := ret
	bresp, _ := json.Marshal(obj)

	fmt.Fprintf(httpresponsewriter, string(bresp)) // write data to response

}

// ForgotPasswordPage is for the user to signup
func ForgotPasswordPage(httpresponsewriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, sysid string) {

	type ControllerInfo struct {
		Name    string
		Message string
	}
	type DisplayTemplate struct {
		Info ControllerInfo
	}

	items := DisplayTemplate{}
	items.Info.Name = "Forgot Password"

	if req.Method != "POST" {

		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = ""
		t.Execute(httpresponsewriter, items)

		return
	}

	// Check if user wants a code or has submited the code.
	//
	// buttonpressed := req.FormValue("submit")      // will be "submit1" or "submit2"

	// applicationid := req.FormValue("applicationid")
	applicationid := "FestaJunina" // antes "restaurante"

	usernamemix := req.FormValue("username") // email address
	requestcode := req.FormValue("requestcode")
	password := req.FormValue("password")
	passwordvalidate := req.FormValue("passwordvalidate")
	buttonpressed := req.FormValue("requestcode") // will be "submit1" or "submit2"

	if buttonpressed == "requestcode" {
		log.Println(usernamemix + " > Code: 123456")
	} else {
		log.Println(usernamemix + " > Update password")
	}

	username := strings.ToUpper(usernamemix)

	if username == "" {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Please_enter_Name."
		t.Execute(httpresponsewriter, items)
		return
	}

	if requestcode == "" {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Code not supplied."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == "" {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Password_Is_Empty."
		t.Execute(httpresponsewriter, items)
		return
	}

	if passwordvalidate == "" {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Password_Validate_is_Empty."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == passwordvalidate {

		// Call API to check if user exists and create
		var resultado = ResetPassword(sysid, redisclient, username, requestcode, password, passwordvalidate, applicationid)
		if resultado.ErrorCode == "200 OK" {

		} else {
			t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
			items.Info.Message = "Passwords_mismatch."
			t.Execute(httpresponsewriter, items)
			return
		}

		http.Redirect(httpresponsewriter, req, "/", 303)
	} else {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Passwords_mismatch."
		t.Execute(httpresponsewriter, items)
		return
	}

}

// UpdatePasswordForgotten is for the user to signup
func UpdatePasswordForgotten(httpresponsewriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, sysid string) {

	type ControllerInfo struct {
		Name         string
		Message      string
		EmailAddress string
	}
	type DisplayTemplate struct {
		Info ControllerInfo
	}

	items := DisplayTemplate{}
	items.Info.Name = "Forgot Password"

	applicationid := "FestaJunina" // antes "restaurante"

	usernamemix := req.FormValue("username") // email address
	items.Info.EmailAddress = "test@email.com"

	requestcode := req.FormValue("requestcode")
	password := req.FormValue("password")
	passwordvalidate := req.FormValue("passwordvalidate")

	log.Println(usernamemix + " > Code: 123456")

	username := strings.ToUpper(usernamemix)

	if username == "" {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Please_enter_Name."
		t.Execute(httpresponsewriter, items)
		return
	}

	if requestcode == "" {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Code not supplied."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == "" {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Password_Is_Empty."
		t.Execute(httpresponsewriter, items)
		return
	}

	if passwordvalidate == "" {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Password_Validate_is_Empty."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == passwordvalidate {

		// Call API to check if user exists and create
		var resultado = ResetPassword(sysid, redisclient, username, requestcode, password, passwordvalidate, applicationid)
		if resultado.ErrorCode == "200 OK" {

		} else {
			t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
			items.Info.Message = "Passwords_mismatch."
			t.Execute(httpresponsewriter, items)
			return
		}

		http.Redirect(httpresponsewriter, req, "/", 303)
	} else {
		t, _ := template.ParseFiles("templates/security/forgotpasswordheader.html", "templates/security/forgotpassworddetail.html")
		items.Info.Message = "Passwords_mismatch."
		t.Execute(httpresponsewriter, items)
		return
	}

}

// RequestCode is for the user to signup
func RequestCode(httpresponsewriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, sysid string) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	GetResetPasswordCode(sysid, redisclient, bodybyte)

}

// ChangePassword is for the user to change password
func ChangePassword(httpresponsewriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, sysid string) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	ret := ApplyNewPassword(sysid, redisclient, bodybyte)

	if ret.ErrorDescription == "Code has expired." {
		obj := ret
		bresp, _ := json.Marshal(obj)

		fmt.Fprintf(httpresponsewriter, string(bresp)) // write data to response
	} else {

		obj := ret
		bresp, _ := json.Marshal(obj)

		fmt.Fprintf(httpresponsewriter, string(bresp)) // write data to response
	}
}

// UserRolesUpdate is for the user to update roles
func UserRolesUpdate(httpresponsewriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, sysid string) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	ret := UpdateUserRoles(sysid, redisclient, bodybyte)

	if ret.ErrorDescription == "Code has expired." {
		obj := ret
		bresp, _ := json.Marshal(obj)

		fmt.Fprintf(httpresponsewriter, string(bresp)) // write data to response
	} else {

		obj := ret
		bresp, _ := json.Marshal(obj)

		fmt.Fprintf(httpresponsewriter, string(bresp)) // write data to response
	}
}

// LogoutPage is for the user to logout
func LogoutPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	// 2018-Sep-02
	// Eu tenho que chamar o server para fazer logout e deletar o entry to redis e db se tiver
	// e nao apenas limpar o cookie

	cookie, _ := httprequest.Cookie("DanBTCjwt")
	if cookie != nil {
		if cookie.Value != "Anonymous" {
			c := &http.Cookie{
				Name:     "DanBTCjwt",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: true,
			}
			http.SetCookie(httpresponsewriter, c)
		}
	}

	http.Redirect(httpresponsewriter, httprequest, "/", 303)
}

// LoginPage is for login
func LoginPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client, sysid string) {

	type ControllerInfo struct {
		Name    string
		Message string
	}
	type DisplayTemplate struct {
		Info ControllerInfo
	}

	items := DisplayTemplate{}
	items.Info.Name = "Login Page"

	if httprequest.Method != "POST" {

		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = ""
		t.Execute(httpresponsewriter, items)

		// http.ServeFile(httpresponsewriter, httprequest, "templates/security/login.html")
		return
	}

	usernamemix := httprequest.FormValue("userid")
	password := httprequest.FormValue("password")

	userid := strings.ToUpper(usernamemix)

	if userid == "" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Enter email address and password."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == "" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Enter email address and password."
		t.Execute(httpresponsewriter, items)
		return
	}

	cookiekeyJWT := "DanBTCjwt"
	cookiekeyUSERID := "DanBTCuserid"

	cookieJWT, _ := httprequest.Cookie(cookiekeyJWT)
	cookieUSERID, _ := httprequest.Cookie(cookiekeyUSERID)

	if cookieJWT != nil {
		cokJWT := &http.Cookie{
			Name:     cookiekeyJWT,
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, cokJWT)
	}

	if cookieUSERID != nil {
		cokUSERID := &http.Cookie{
			Name:     cookiekeyUSERID,
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, cokUSERID)
	}

	// Check if the user is valid and issue reference token
	//
	var resultado = LoginUserV2(sysid, redisclient, userid, password)

	if resultado.JWT == "Error" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Login error. Try again."
		t.Execute(httpresponsewriter, items)
		return
	}

	if resultado.ApplicationID != "Restaurante" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "User is invalid."
		t.Execute(httpresponsewriter, items)
		return
	}

	// Store Token in Cache
	var jwttoken = resultado.JWT
	year, month, day := time.Now().Date()
	var expiry = strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + strconv.Itoa(int(day))

	rediskey := "DanBTCjwt" + userid

	var credentials commonstruct.Credentials
	credentials.UserID = userid
	credentials.KeyJWT = rediskey
	credentials.JWT = jwttoken
	credentials.Expiry = expiry
	credentials.ClaimSet = resultado.ClaimSet
	credentials.ApplicationID = resultado.ApplicationID
	credentials.Name = resultado.Name
	credentials.IsAdmin = resultado.IsAdmin
	credentials.CentroID = resultado.CentroID

	jsonval, _ := json.Marshal(credentials)
	jsonstring := string(jsonval)

	_ = redisclient.Set(rediskey, jsonstring, 0).Err()

	// store in cookie
	// 2 hours ==> 4 hours
	expiration := time.Now().Add(4 * time.Hour)

	cokJWT := &http.Cookie{
		Name:     cookiekeyJWT,
		Value:    jwttoken,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, cokJWT)

	cokUSERID := &http.Cookie{
		Name:     cookiekeyUSERID,
		Value:    userid,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, cokUSERID)

	http.Redirect(httpresponsewriter, httprequest, "/", 303)

	return
}

// AnonymousLogin is for login
func AnonymousLogin(httpresponsewriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client, useridin string, username string) {
	log.Println("AnonymousLogin Called " + useridin)

	userid := strings.ToUpper(useridin)
	log.Println("AnonymousLogin - User ID: " + userid)

	cookiekeyJWT := "DanBTCjwt"
	cookiekeyUSERID := "DanBTCuserid"
	cookiealreadystored := "No"

	cookieJWT, _ := httprequest.Cookie(cookiekeyJWT)
	cookieUSERID, _ := httprequest.Cookie(cookiekeyUSERID)

	if cookieJWT != nil {
		//-------------------------------------------------------
		// Neste caso apenas retorne o cookie value, nao apague
		//-------------------------------------------------------

		// ??? Este e' o cookie que armazena a JWT
	}

	if cookieUSERID != nil {
		//-------------------------------------------------------
		// Neste caso apenas retorne o cookie value, nao apague
		//-------------------------------------------------------

		userid = cookieUSERID.Value
		cookiealreadystored = "Yes"
	}

	resultado := commonstruct.Credentials{}
	resultado.JWT = "Anonymous"
	resultado.ApplicationID = "Restaurante"

	// Store Token in Cache
	var jwttoken = resultado.JWT
	year, month, day := time.Now().Date()
	var expiry = strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + strconv.Itoa(int(day))

	rediskey := "DanBTCjwt" + userid

	var credentials commonstruct.Credentials
	credentials.UserID = userid
	credentials.KeyJWT = rediskey
	credentials.JWT = jwttoken
	credentials.Expiry = expiry
	credentials.ClaimSet = resultado.ClaimSet
	credentials.ApplicationID = resultado.ApplicationID
	credentials.Name = username
	credentials.IsAdmin = resultado.IsAdmin
	credentials.CentroID = resultado.CentroID

	jsonval, _ := json.Marshal(credentials)
	jsonstring := string(jsonval)

	// ---------------------------------------
	//         Store in cache
	// ---------------------------------------
	_ = redisclient.Set(rediskey, jsonstring, 0).Err()

	if cookiealreadystored == "No" {

		// store in cookie
		// 1 month
		expiration := time.Now().Add(720 * time.Hour)
		// expiration := time.Now().Add(1 * time.Hour)

		cokJWT := &http.Cookie{
			Name:     cookiekeyJWT,
			Value:    jwttoken,
			Path:     "/",
			Expires:  expiration,
			MaxAge:   0,
			HttpOnly: true,
		}

		http.SetCookie(httpresponsewriter, cokJWT)
		log.Println("Storing Cookie: " + cookiekeyJWT)

		cokUSERID := &http.Cookie{
			Name:     cookiekeyUSERID,
			Value:    userid,
			Path:     "/",
			Expires:  expiration,
			MaxAge:   0,
			HttpOnly: true,
		}

		http.SetCookie(httpresponsewriter, cokUSERID)
		log.Println("Storing Cookie: " + cookiekeyUSERID)
	} else {

		log.Println("Reusing Cookie ! ")

	}

	// http.Redirect(httpresponsewriter, httprequest, "/", 303)

	return
}

// ControllerInfo is
type ControllerInfo struct {
	Name          string
	Message       string
	UserID        string
	UserName      string
	ApplicationID string //
	IsAdmin       string //
}

// Row is
type Row struct {
	Description []string
}

// DisplayTemplate is
type DisplayTemplate struct {
	Info       ControllerInfo
	FieldNames []string
	Rows       []Row
	Pratos     []dishes.Dish
}

// Instructions is for login
func Instructions(httpresponsewriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

	// create new template
	t, error := template.ParseFiles("html/homepage.html", "templates/main/instructions.html")

	if error != nil {
		panic(error)
	}

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Instructions"

	t.Execute(httpresponsewriter, items)

}
