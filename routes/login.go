package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Login shows the login page.
*/
func Login(c *gin.Context) {
	if auth.GetUserMode(c) == "guest" {
		HeaderData := &web.HeaderData{Title: "Login", StyleSheets: []string{"global"}}
		loggedIn := true
		c.HTML(http.StatusOK, "login.tmpl", gin.H{"loggedIn": loggedIn, "HeaderData": HeaderData})
	} else {
		c.Redirect(http.StatusSeeOther, "/dashboard")
	}
}

/*
LoginPOST logs a user in.
*/
func LoginPOST(c *gin.Context) {
	HeaderData := &web.HeaderData{Title: "Login", StyleSheets: []string{"global"}}
	c.Request.ParseForm()
	username := c.PostForm("username")
	password := c.PostForm("password")
	//The login function returns the uuid but returns a blank string if it fails
	loggedIn, _ := db.UserLogin(username, password)
	if loggedIn {
		userData, _ := db.UserQuery(username)
		auth.SetLogin(c, userData.UserID)
		c.Redirect(http.StatusSeeOther, "/dashboard") // Although gin's method here is named Redirect, the HTTP response code used is 303. See https://en.wikipedia.org/wiki/HTTP_303 for more information.
	} else {
		c.HTML(200, "login.tmpl", gin.H{"loggedIn": false, "HeaderData": HeaderData})
	}
}

//Logout logs a user out by voiding the login cookie
func Logout(c *gin.Context) {
	auth.SetLogin(c, "")
	c.Redirect(http.StatusSeeOther, "/login")
}
