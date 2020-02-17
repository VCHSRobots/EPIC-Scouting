package routes

import (
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Login shows the login page.
*/
func Login(c *gin.Context) {
	HeaderData := &web.HeaderData{Title: "Login", StyleSheets: []string{"global"}}
	loggedIn := true
	c.HTML(http.StatusOK, "login.tmpl", gin.H{"loggedIn": loggedIn, "HeaderData": HeaderData})
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
		uuid := db.GetUserID(username)
		cookie := fmt.Sprintf("%s %s", uuid, username)
		//secure cannot be set to true until we get http working
		//TODO: put in actual site domain
		c.SetCookie("login", cookie, 3600, "/", "", http.SameSiteLaxMode, false, false)
		c.Redirect(http.StatusSeeOther, "/dashboard") // Although gin's method here is named Redirect, the HTTP response code used is 303. See https://en.wikipedia.org/wiki/HTTP_303 for more information.
	} else {
		c.HTML(200, "login.tmpl", gin.H{"loggedIn": false, "HeaderData": HeaderData})
	}
}

//Logout logs a user out by voiding the login cookie
func Logout(c *gin.Context) {
	c.SetCookie("login", "", 1, "/", "", http.SameSiteLaxMode, false, false)
	c.Redirect(http.StatusSeeOther, "/login")
}
