// Package activitieshandler API calls for dishes web
// --------------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/dishesapicalls.go
// --------------------------------------------------------------
package activitieshandler

import (
	"encoding/json"
	"feijoadajusu/areas/commonstruct"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	helper "feijoadajusu/areas/helper"
	activities "feijoadajusu/models"
)

// List works
func actlist() []activities.Activity {

	var apiserver string
	apiserver = helper.Getvaluefromcache("MSAPIactivitiesIPAddress")

	urlrequest := apiserver + "/list"

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []activities.Activity

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

	var dishlist []activities.Activity

	if err := json.NewDecoder(resp.Body).Decode(&dishlist); err != nil {
		log.Println(err)
	}

	return dishlist
}

// ListV2 sending error back
func actlistV2() ([]activities.Activity, commonstruct.Resultado) {

	var resultado commonstruct.Resultado
	resultado.ErrorCode = "0001"
	resultado.ErrorDescription = "Successful transaction"
	resultado.IsSuccessful = "true"

	var apiserver string

	apiserver = helper.Getvaluefromcache("MSAPIactivitiesIPAddress")

	urlrequest := apiserver + "/list"

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []activities.Activity

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

	var dishlist []activities.Activity

	if err := json.NewDecoder(resp.Body).Decode(&dishlist); err != nil {
		log.Println(err)
	}

	resultado.ErrorCode = "0001"
	resultado.ErrorDescription = "Successful transaction"
	resultado.IsSuccessful = "true"
	return dishlist, resultado
}

// APIcallAdd is
func APIcallAdd(objInsert activities.Activity) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	// mongodbvar.APIServer, _ = redisclient.Get(sysid + "MSAPIactivitiesIPAddress").Result()
	mongodbvar.APIServer = helper.Getvaluefromcache("MSAPIactivitiesIPAddress")

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := mongodbvar.APIServer
	resource := "/add"

	data := url.Values{}
	data.Add("name", objInsert.Name)
	data.Add("type", objInsert.Type)
	data.Add("status", objInsert.Status)
	data.Add("description", objInsert.Description)
	data.Add("startdate", objInsert.StartDate)
	data.Add("enddate", objInsert.EndDate)

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
func FindAPI(objFind string) activities.Activity {

	var apiserver string
	// apiserver, _ = redisclient.Get(sysid + "MSAPIactivitiesIPAddress").Result()
	apiserver = helper.Getvaluefromcache("MSAPIactivitiesIPAddress")

	// This is essential! Because if the string has spaces it doesn't work without the escape
	// Bolo de Cenoura = Bolo+de+Cenoura   >>> Works as a dream!
	objfindescaped := url.QueryEscape(objFind)
	urlrequest := apiserver + "/find?name=" + objfindescaped

	urlrequestencoded, _ := url.ParseRequestURI(urlrequest)
	// url := fmt.Sprintf(urlrequest)
	url := urlrequestencoded.String()
	// tw.Text = strings.Replace(tw.Text, " ", "+", -1)
	// urlx := url.QueryEscape(urlrequest)

	var emptydisplay activities.Activity

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

	var activitiesback activities.Activity

	if err := json.NewDecoder(resp.Body).Decode(&activitiesback); err != nil {
		log.Println(err)
	}

	return activitiesback

}

// FindActiveAPI is to find the active activity/ event
func FindActiveAPI() activities.Activity {

	var apiserver string
	// apiserver, _ = redisclient.Get(sysid + "MSAPIactivitiesIPAddress").Result()
	apiserver = helper.Getvaluefromcache("MSAPIactivitiesIPAddress")

	urlrequest := apiserver + "/findactive"

	urlrequestencoded, _ := url.ParseRequestURI(urlrequest)
	url := urlrequestencoded.String()

	var emptydisplay activities.Activity

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

		log.Println("FindActiveAPI Error client.Do(req): ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	var activitiesback activities.Activity

	if err := json.NewDecoder(resp.Body).Decode(&activitiesback); err != nil {
		log.Println(err)
	}

	return activitiesback

}

// UpdateAPI is
func UpdateAPI(objUpdate activities.Activity) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	mongodbvar.APIServer = helper.Getvaluefromcache("MSAPIactivitiesIPAddress")

	apiURL := mongodbvar.APIServer
	resource := "/update"

	data := url.Values{}
	data.Add("name", objUpdate.Name)
	data.Add("type", objUpdate.Type)
	data.Add("status", objUpdate.Status)
	data.Add("description", objUpdate.Description)
	data.Add("startdate", objUpdate.StartDate)
	data.Add("enddate", objUpdate.EndDate)

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

// DeleteAPI is
func DeleteAPI(objUpdate activities.Activity) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	mongodbvar.APIServer = helper.Getvaluefromcache("MSAPIactivitiesIPAddress")

	apiURL := mongodbvar.APIServer
	resource := "/delete"

	data := url.Values{}
	data.Add("name", objUpdate.Name)

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

// DeleteMultipleAPI is
func DeleteMultipleAPI(objtodelete []string) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	mongodbvar.APIServer = helper.Getvaluefromcache("MSAPIactivitiesIPAddress")

	apiURL := mongodbvar.APIServer
	resource := "/delete"

	data := url.Values{}
	data.Add("name", objtodelete[0])

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

// ListActivities works
func ListActivities() []activities.Activity {

	var apiserver string
	apiserver = helper.Getvaluefromcache("MSAPIactivitiesIPAddress")

	urlrequest := apiserver + "/list"

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []activities.Activity

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

	var list []activities.Activity

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}
