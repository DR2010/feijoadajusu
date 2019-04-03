// Package helper Handler for dishes web
// -----------------------------------------------------------
// .../src/restauranteweb/areas/helper/roothandler.go
// -----------------------------------------------------------
package helper

import (
	"feijoadajusu/areas/commonstruct"
	"feijoadajusu/models"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
)

// This is the template to display as part of the html template
//

// ControllerInfo is
type ControllerInfo struct {
	UserID       string
	Name         string
	Message      string
	Currency     string
	FromDate     string
	ToDate       string
	Application  string
	IsAdmin      string //
	Organisation string //
	Database     string //
}

// Row is
type Row struct {
	Description []string
}

// DisplayTemplate is
type DisplayTemplate struct {
	Info ControllerInfo
}

var mongodbvar commonstruct.DatabaseX

// HomePage = assemble results of API call to dish list
//
func HomePage(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials models.Credentials) {

	// create new template
	t, _ := template.ParseFiles("html/homepage.html", "templates/main/pagebodytemplate.html")

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = credentials.Name
	items.Info.UserID = credentials.UserID
	items.Info.Application = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin
	items.Info.Organisation = Getvaluefromcache("Organisation")

	fmt.Println("Organisation: " + items.Info.Organisation)

	org := Getvaluefromcache("Organisation")
	items.Info.Organisation = org

	db := Getvaluefromcache("APIMongoDBDatabase")
	items.Info.Database = db

	t.Execute(httpwriter, items)
}

// HomePage2 = assemble results of API call to dish list
//
func HomePage2(httpwriter http.ResponseWriter) {

	// create new template
	var listtemplate = `
			{{define "listtemplate"}}
			This is my web site, Daniel - aka D#.
			<p/>
			<p/>
			<picture>
				<img src="images/avatar.png" alt="Avatar" width="400" height="400">
			</picture>
			{{end}}
			`

	t, _ := template.ParseFiles("html/index.html")
	t, _ = t.Parse(listtemplate)

	t.Execute(httpwriter, listtemplate)
}
