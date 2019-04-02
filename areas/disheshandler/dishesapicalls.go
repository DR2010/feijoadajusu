// Package disheshandler API calls for dishes web
// --------------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/dishesapicalls.go
// --------------------------------------------------------------
package disheshandler

import (
	"encoding/json"
	"feijoadajusu/areas/commonstruct"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	helper "feijoadajusu/areas/helper"
	dishes "feijoadajusu/models"
)

// // Dish is to be exported
// type Dish struct {
// 	SystemID         bson.ObjectId `json:"id"        bson:"_id,omitempty"`
// 	Name             string        // name of the dish - this is the KEY, must be unique
// 	Type             string        // type of dish, includes drinks and deserts
// 	Price            string        // preco do prato multiplicar por 100 e nao ter digits
// 	GlutenFree       string        // Gluten free dishes
// 	DairyFree        string        // Dairy Free dishes
// 	Vegetarian       string        // Vegeterian dishes
// 	InitialAvailable string        // Number of items initially available
// 	CurrentAvailable string        // Currently available
// 	ImageName        string        // Image Name
// }

// ListDishes works
func listdishes() []dishes.Dish {

	var apiserver string
	// apiserver, _ = redisclient.Get("MSAPIdishesIPAddress").Result()
	apiserver = helper.Getvaluefromcache("MSAPIdishesIPAddress")

	urlrequest := apiserver + "/dishlist"

	// urlrequest = "http://localhost:1520/dishlist"

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []dishes.Dish

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		// log.Fatal("Do: ", err)
		fmt.Println("End point not availabe: " + urlrequest + " - Error: " + err.Error())
		return emptydisplay
	}

	defer resp.Body.Close()

	var dishlist []dishes.Dish

	if err := json.NewDecoder(resp.Body).Decode(&dishlist); err != nil {
		log.Println(err)
	}

	return dishlist
}

// ListDishes sending error back
func listdishesV2() ([]dishes.Dish, commonstruct.Resultado) {

	var resultado commonstruct.Resultado
	resultado.ErrorCode = "0001"
	resultado.ErrorDescription = "Successful transaction"
	resultado.IsSuccessful = "true"

	var apiserver string
	// apiserver, _ = redisclient.Get(sysid + "MSAPIdishesIPAddress").Result()

	apiserver = helper.Getvaluefromcache("MSAPIdishesIPAddress")

	urlrequest := apiserver + "/dishlist"

	// urlrequest = "http://localhost:1520/dishlist"

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []dishes.Dish

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay, resultado
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		// log.Fatal("Do: ", err)
		log.Println(err)
		resultado.ErrorDescription = " = Error: Dishes List not available, please try later."
		resultado.IsSuccessful = "false"
		resultado.ErrorCode = "0102" // can't reach destination
		return emptydisplay, resultado
	}

	defer resp.Body.Close()

	var dishlist []dishes.Dish

	if err := json.NewDecoder(resp.Body).Decode(&dishlist); err != nil {
		log.Println(err)
	}

	resultado.ErrorCode = "0001"
	resultado.ErrorDescription = "Successful transaction"
	resultado.IsSuccessful = "true"
	return dishlist, resultado
}

// APIcallAdd is
func APIcallAdd(dishInsert dishes.Dish) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	// mongodbvar.APIServer, _ = redisclient.Get(sysid + "MSAPIdishesIPAddress").Result()
	mongodbvar.APIServer = helper.Getvaluefromcache("MSAPIdishesIPAddress")

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := mongodbvar.APIServer
	resource := "/dishadd"

	data := url.Values{}
	data.Add("dishname", dishInsert.Name)
	data.Add("dishtype", dishInsert.Type)
	data.Add("dishprice", dishInsert.Price)
	data.Add("dishglutenfree", dishInsert.GlutenFree)
	data.Add("dishdairyfree", dishInsert.DairyFree)
	data.Add("dishvegetarian", dishInsert.Vegetarian)
	data.Add("dishinitialavailable", dishInsert.InitialAvailable)
	data.Add("dishcurrentavailable", dishInsert.CurrentAvailable)
	data.Add("dishimagename", dishInsert.ImageName)
	data.Add("dishdescription", dishInsert.Description)
	data.Add("dishdescricao", dishInsert.Descricao)
	data.Add("dishactivitytype", dishInsert.ActivityType)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)
	fmt.Println("body:" + data.Encode())

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	} else {
		emptydisplay.IsSuccessful = "N"
	}

	return emptydisplay

}

// FindAPI is to find stuff
func FindAPI(dishFind string) dishes.Dish {

	var apiserver string
	// apiserver, _ = redisclient.Get(sysid + "MSAPIdishesIPAddress").Result()
	apiserver = helper.Getvaluefromcache("MSAPIdishesIPAddress")

	// This is essential! Because if the string has spaces it doesn't work without the escape
	// Bolo de Cenoura = Bolo+de+Cenoura   >>> Works as a dream!
	dishfindescaped := url.QueryEscape(dishFind)
	urlrequest := apiserver + "/dishfind?dishname=" + dishfindescaped

	urlrequestencoded, _ := url.ParseRequestURI(urlrequest)
	// url := fmt.Sprintf(urlrequest)
	url := urlrequestencoded.String()
	// tw.Text = strings.Replace(tw.Text, " ", "+", -1)
	// urlx := url.QueryEscape(urlrequest)

	var emptydisplay dishes.Dish

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// log.Fatal("NewRequest: ", err)
		log.Println("FindAPI Error http.NewRequest(GET, url, nil): ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		// log.Fatal("Do: ", err)

		log.Println("FindAPI Error client.Do(req): ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	var dishback dishes.Dish

	if err := json.NewDecoder(resp.Body).Decode(&dishback); err != nil {
		log.Println(err)
	}

	return dishback

}

// DishupdateAPI is
func DishupdateAPI(dishUpdate dishes.Dish) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	// mongodbvar.APIServer, _ = redisclient.Get(sysid + "MSAPIdishesIPAddress").Result()
	mongodbvar.APIServer = helper.Getvaluefromcache("MSAPIdishesIPAddress")

	apiURL := mongodbvar.APIServer
	resource := "/dishupdate"

	data := url.Values{}
	data.Add("dishname", dishUpdate.Name)
	data.Add("dishtype", dishUpdate.Type)
	data.Add("dishprice", dishUpdate.Price)
	data.Add("dishglutenfree", dishUpdate.GlutenFree)
	data.Add("dishdairyfree", dishUpdate.DairyFree)
	data.Add("dishvegetarian", dishUpdate.Vegetarian)
	data.Add("dishinitialavailable", dishUpdate.InitialAvailable)
	data.Add("dishcurrentavailable", dishUpdate.CurrentAvailable)
	data.Add("dishimagename", dishUpdate.ImageName)
	data.Add("dishdescription", dishUpdate.Description)
	data.Add("dishdescricao", dishUpdate.Descricao)
	data.Add("dishactivitytype", dishUpdate.ActivityType)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	return emptydisplay
}

// DishdeleteAPI is
func DishdeleteAPI(dishUpdate dishes.Dish) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	// mongodbvar.APIServer, _ = redisclient.Get(sysid + "MSAPIdishesIPAddress").Result()
	mongodbvar.APIServer = helper.Getvaluefromcache("MSAPIdishesIPAddress")

	apiURL := mongodbvar.APIServer
	resource := "/dishdelete"

	data := url.Values{}
	data.Add("dishname", dishUpdate.Name)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	return emptydisplay
}

// DishDeleteMultipleAPI is
func DishDeleteMultipleAPI(dishestodelete []string) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	// mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	// mongodbvar.APIServer, _ = redisclient.Get(sysid + "MSAPIdishesIPAddress").Result()
	mongodbvar.APIServer = helper.Getvaluefromcache("MSAPIdishesIPAddress")

	apiURL := mongodbvar.APIServer
	resource := "/dishdelete"

	data := url.Values{}
	data.Add("dishname", dishestodelete[0])

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	return emptydisplay
}

// Listdishes works
func Listdishes() []dishes.Dish {

	var apiserver string
	// apiserver, _ = redisclient.Get(sysid + "MSAPIdishesIPAddress").Result()
	apiserver = helper.Getvaluefromcache("MSAPIdishesIPAddress")

	urlrequest := apiserver + "/dishlist"

	// urlrequest = "http://localhost:1520/dishlist"

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []dishes.Dish

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

	var dishlist []dishes.Dish

	if err := json.NewDecoder(resp.Body).Decode(&dishlist); err != nil {
		log.Println(err)
	}

	return dishlist
}
