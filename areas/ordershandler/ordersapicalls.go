// Package ordershandler API calls for dishes web
// --------------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/orderapicalls.go
// --------------------------------------------------------------
package ordershandler

import (
	"encoding/json"
	"feijoadajusu/areas/commonstruct"
	"feijoadajusu/models"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2/bson"
)

// Dish is to be exported
type Dish struct {
	SystemID   bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Name       string        // name of the dish - this is the KEY, must be unique
	Type       string        // type of dish, includes drinks and deserts
	Price      string        // preco do prato multiplicar por 100 e nao ter digits
	GlutenFree string        // Gluten free dishes
	DairyFree  string        // Dairy Free dishes
	Vegetarian string        // Vegeterian dishes
}

// SearchCriteria is what the client wants
type SearchCriteria struct {
	ID                   string // random ID for order, yet to define algorithm
	ClientName           string // Client Name
	ClientID             string // Client ID in case they logon
	Date                 string // Order Date
	Time                 string // Order Time
	Status               string // Open, Completed, Cancelled
	EatMode              string // EatIn, TakeAway, Delivery
	DeliveryMode         string // Internal, UberEats,
	DeliveryFee          string // Delivery Fee
	DeliveryLocation     string // Address
	DeliveryContactPhone string // Delivery phone number
}

// RespAddOrder I am not sure
type RespAddOrder struct {
	ID       string
	ClientID string
}

// FindAPI is to find stuff
func FindAPI(sysid string, redisclient *redis.Client, orderFind string) models.Order {

	var apiserver string
	apiserver, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()
	urlrequest := apiserver + "/orderfind?orderid=" + orderFind

	url := fmt.Sprintf(urlrequest)

	var emptydisplay models.Order

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	var orderback models.Order

	if err := json.NewDecoder(resp.Body).Decode(&orderback); err != nil {
		log.Println(err)
	}

	return orderback

}

// APICallList works
// Order List
func APICallList(sysid string, redisclient *redis.Client) []models.Order {

	var apiserver string
	var emptydisplay []models.Order

	apiserver, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()
	urlrequest := apiserver + "/orderlist"

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	// return list of orders
	var list []models.Order

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// APICallListV2 works
// Order List
func APICallListV2(sysid string, redisclient *redis.Client, credentials models.Credentials) []models.Order {

	var apiserver string
	var emptydisplay []models.Order

	apiserver, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()

	urlrequest := apiserver + "/orderlist?clientid=" + credentials.UserID

	// Check if user is admin
	for x := 0; x < len(credentials.ClaimSet); x++ {
		if credentials.ClaimSet[x].Type == "USERTYPE" {
			if credentials.ClaimSet[x].Value == "ADMIN" {
				// list all if user is admin
				urlrequest = apiserver + "/orderlist"
				break
			}
		}
	}

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	// return list of orders
	var list []models.Order

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// APISaveOrderToMySQL works
func APISaveOrderToMySQL(sysid string, redisclient *redis.Client) {

	var apiserver string

	apiserver, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()
	urlrequest := apiserver + "/savetomysql"

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()

	// return list of orders
	var list []models.Order

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return
}

// APICallListCompleted is completed
func APICallListCompleted(sysid string, redisclient *redis.Client, credentials models.Credentials) []models.Order {

	var apiserver string
	var emptydisplay []models.Order

	apiserver, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()
	urlrequest := apiserver + "/ordercompleted"

	log.Println("APICallListCompleted " + urlrequest)

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	// return list of orders
	var list []models.Order

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// APICallListStatus is completed
func APICallListStatus(sysid string, redisclient *redis.Client, credentials models.Credentials, status string) []models.Order {

	var apiserver string
	var emptydisplay []models.Order

	apiserver, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()

	urlrequest := apiserver + "/orderstatus?status=" + status

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	// return list of orders
	var list []models.Order

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// APICallListStatusActivity is completed
func APICallListStatusActivity(sysid string, redisclient *redis.Client, credentials models.Credentials, status string, activity string) []models.Order {

	var apiserver string
	var emptydisplay []models.Order

	apiserver, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()

	urlrequest := apiserver + "/orderstatusactivity?status=" + status + "&activity=" + activity

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	// return list of orders
	var list []models.Order

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// APICallAdd is
func APICallAdd(sysid string, redisclient *redis.Client, bodybyte []byte) RespAddOrder {

	envirvar := new(commonstruct.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/orderadd"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(bodystr)
	resp2, err := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()
	var objectback RespAddOrder

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
		var resultado = resp2.Body
		log.Println(resultado)

		if err = json.NewDecoder(resp2.Body).Decode(&objectback); err != nil {
			log.Println(err)
		} else {

			var x = objectback.ID
			log.Println(x)
		}

	} else {
		emptydisplay.IsSuccessful = "N"

	}
	return objectback
}

// APICallAddOrderClient is designed to add an order and an anonymous client IF required (id not passed in)!
func APICallAddOrderClient(sysid string, redisclient *redis.Client, bodybyte []byte) RespAddOrder {

	envirvar := new(commonstruct.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/APIorderadd"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(bodystr)
	resp2, err := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()
	var objectback RespAddOrder

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
		var resultado = resp2.Body
		log.Println(resultado)

		if err = json.NewDecoder(resp2.Body).Decode(&objectback); err != nil {
			log.Println(err)
		} else {

			var orderid = objectback.ID
			var clientid = objectback.ClientID
			log.Println("Order ID: " + orderid)
			log.Println("Client ID: " + clientid)

		}

	} else {
		emptydisplay.IsSuccessful = "N"

	}
	return objectback
}

// APICallFind is to find stuff
func APICallFind(sysid string, redisclient *redis.Client, objectfind string) models.Order {

	var apiserver string
	apiserver, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()

	urlrequest := apiserver + "/orderfind?ID=" + objectfind

	url := fmt.Sprintf(urlrequest)

	var emptydisplay models.Order

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	var objectback models.Order

	if err := json.NewDecoder(resp.Body).Decode(&objectback); err != nil {
		log.Println(err)
	}

	return objectback

}

// APICallUpdate is
func APICallUpdate(sysid string, redisclient *redis.Client, bodybyte []byte) RespAddOrder {

	envirvar := new(commonstruct.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get(sysid + "MSAPIordersIPAddress").Result()

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/orderupdate"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(bodystr)
	resp2, err := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()
	var objectback RespAddOrder

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
		var resultado = resp2.Body
		log.Println(resultado)

		if err = json.NewDecoder(resp2.Body).Decode(&objectback); err != nil {
			log.Println(err)
		} else {

			var x = objectback.ID
			log.Println(x)
		}

	} else {
		emptydisplay.IsSuccessful = "N"

	}
	return objectback
}
