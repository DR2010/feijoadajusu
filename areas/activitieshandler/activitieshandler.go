// Package activitieshandler Handler for dishes web
// -----------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/disheshandler.go
// -----------------------------------------------------------
package activitieshandler

import (
	"feijoadajusu/areas/commonstruct"
	models "feijoadajusu/models"
	"html/template"
	"mongodb/dishes"
	"net/http"

	"github.com/go-redis/redis"
)

// This is the template to display as part of the html template
//

// ControllerInfo is
type ControllerInfo struct {
	Name          string
	Message       string
	UserID        string
	UserName      string
	ApplicationID string //
	IsAdmin       string //
	IsAnonymous   string //
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

var mongodbvar commonstruct.DatabaseX

// List = assemble results of API call to dish list
//
func List(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials models.Credentials, sysid string) {

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/activity/listtemplate.html")

	// Get list of dishes (api call)
	//
	// var dishlist = listdishes(redisclient)
	actlist, error := actlistV2()

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Activity List"
	items.Info.UserID = credentials.UserID
	items.Info.UserName = credentials.Name
	items.Info.ApplicationID = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin

	if error.IsSuccessful == "false" {

		items.Info.Name = "Activity List " + error.ErrorDescription

		// do something
	}

	var numberoffields = 6

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Name"
	items.FieldNames[1] = "Type"
	items.FieldNames[2] = "Description"
	items.FieldNames[3] = "Status"
	items.FieldNames[4] = "StartDate"
	items.FieldNames[5] = "EndDate"

	// Set rows to be displayed
	items.Rows = make([]Row, len(actlist))
	// items.RowID = make([]int, len(actlist))

	for i := 0; i < len(actlist); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = actlist[i].Name
		items.Rows[i].Description[1] = actlist[i].Type
		items.Rows[i].Description[2] = actlist[i].Description
		items.Rows[i].Description[3] = actlist[i].Status
		items.Rows[i].Description[4] = actlist[i].StartDate
		items.Rows[i].Description[5] = actlist[i].EndDate
	}

	t.Execute(httpwriter, items)
}

// LoadDisplayForAdd is X
func LoadDisplayForAdd(httpwriter http.ResponseWriter) {

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/activity/add.html")

	items := DisplayTemplate{}
	items.Info.Name = "Activity Add"

	t.Execute(httpwriter, items)

}

// Add is
func Add(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, sysid string) {

	objaction := models.Activity{}

	objaction.Type = req.FormValue("activitytype")
	objaction.Name = req.FormValue("activityname") // This is the key, must be unique
	objaction.Description = req.FormValue("activitydescription")
	objaction.Status = req.FormValue("activitystatus")
	objaction.StartDate = req.FormValue("activitystartdate")
	objaction.EndDate = req.FormValue("activityenddate")

	ret := APIcallAdd(objaction)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/activitylist", 301)
	} else {

		// create new template
		// t, _ := template.ParseFiles("html/index.html", "templates/error.html")
		t, _ := template.ParseFiles("html/index.html", "templates/activity/add.html")

		items := DisplayTemplate{}
		items.Info.Name = "Error"
		items.Info.Message = "Activity already exists."

		t.Execute(httpwriter, items)

	}
	return
}

// Update dish sent
func Update(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, sysid string) {

	objaction := models.Activity{}

	objaction.Name = req.FormValue("activityname") // This is the key, must be unique
	objaction.Type = req.FormValue("activitytype")
	objaction.Description = req.FormValue("activitydescription")
	objaction.Status = req.FormValue("activitystatus")
	objaction.StartDate = req.FormValue("activitystartdate")
	objaction.EndDate = req.FormValue("activityenddate")

	ret := UpdateAPI(objaction)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}

// LoadDisplayForUpdate is
func LoadDisplayForUpdate(httpwriter http.ResponseWriter, httprequest *http.Request, credentials models.Credentials) {

	httprequest.ParseForm()

	// Get all selected records
	activityselected := httprequest.Form["activities"]

	var numrecsel = len(activityselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, httprequest, "/activitylist", 301)
		return
	}

	type ControllerInfo struct {
		Name        string
		Message     string
		UserID      string
		Currency    string
		Application string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
		Item       models.Activity
	}

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/activity/update.html")

	items := DisplayTemplate{}
	items.Info.Name = "Activity Add"
	items.Info.Currency = "SUMMARY"
	items.Info.UserID = credentials.UserID
	items.Info.Application = credentials.ApplicationID

	items.Item = models.Activity{}
	items.Item.Name = activityselected[0]

	var activitiesfind = models.Activity{}
	var activitiesname = items.Item.Name

	activitiesfind = FindAPI(activitiesname)
	items.Item = activitiesfind

	t.Execute(httpwriter, items)

	return

}

// LoadDisplayForDelete is
func LoadDisplayForDelete(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client, sysid string) {

	httprequest.ParseForm()

	// Get all selected records
	activityselected := httprequest.Form["activities"]

	var numrecsel = len(activityselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, httprequest, "/activitylist", 301)
		return
	}

	type ControllerInfo struct {
		Name    string
		Message string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
		Item       models.Activity
	}

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/activitydelete.html")

	items := DisplayTemplate{}
	items.Info.Name = "Activity Delete"

	items.Item = models.Activity{}
	items.Item.Name = activityselected[0]

	var activityfind = models.Activity{}
	var activityname = items.Item.Name

	activityfind = FindAPI(activityname)
	items.Item = activityfind

	t.Execute(httpwriter, items)

	return

}

// Delete dish sent
func Delete(redisclient *redis.Client, httpwriter http.ResponseWriter, req *http.Request, sysid string) {

	objaction := models.Activity{}

	objaction.Name = req.FormValue("activityname") // This is the key, must be unique
	objaction.Type = req.FormValue("activitytype")
	objaction.Description = req.FormValue("activitydescription")
	objaction.Status = req.FormValue("activitystatus")
	objaction.StartDate = req.FormValue("activitystartdate")
	objaction.EndDate = req.FormValue("activityenddate")

	ret := DeleteAPI(objaction)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}

func activitydeletedisplay(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, sysid string) {

	req.ParseForm()

	// Get all selected records
	activityselected := req.Form["activities"]

	var numrecsel = len(activityselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}

	type ControllerInfo struct {
		Name string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
		Item       models.Activity
	}

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/activity/delete.html")

	items := DisplayTemplate{}
	items.Info.Name = "Activity Delete"

	items.Item = models.Activity{}
	items.Item.Name = activityselected[0]

	var activityfind = models.Activity{}
	var activityname = items.Item.Name

	activityfind = FindAPI(activityname)
	items.Item = activityfind

	t.Execute(httpwriter, items)

	return

}

// // Dishdeletemultiple is to delete multiple dishes
// func Dishdeletemultiple(httpwriter http.ResponseWriter, req *http.Request) {

// 	req.ParseForm()

// 	// Get all selected records
// 	dishselected := req.Form["dishes"]

// 	var numrecsel = len(dishselected)

// 	if numrecsel <= 0 {
// 		http.Redirect(httpwriter, req, "/dishlist", 301)
// 		return
// 	}

// 	dishtodelete := dishes.Dish{}

// 	ret := commonstruct.Resultado{}

// 	for x := 0; x < len(dishselected); x++ {

// 		dishtodelete.Name = dishselected[x]

// 		ret = Dishdelete(mongodbvar, dishtodelete)
// 	}

// 	if ret.IsSuccessful == "Y" {
// 		// http.ServeFile(httpwriter, req, "success.html")
// 		http.Redirect(httpwriter, req, "/dishlist", 301)
// 		return
// 	}

// 	http.Redirect(httpwriter, req, "/dishlist", 301)
// 	return

// }
