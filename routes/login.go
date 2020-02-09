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
	HeaderData := &web.HeaderData{}
	loggedIn := true
	c.HTML(http.StatusOK, "login.tmpl", gin.H{"loggedIn": loggedIn, "HeaderData": HeaderData})
}

/*
LoginPOST logs a user in.
*/
func LoginPOST(c *gin.Context) {
	HeaderData := &web.HeaderData{}
	c.Request.ParseForm()
	username := c.PostForm("username")
	password := c.PostForm("password")
	loggedIn := db.CheckLoginII(username, password)
	fmt.Println(loggedIn)
	if loggedIn {
		c.Redirect(http.StatusSeeOther, "/dashboard") // Although gin's method here is named Redirect, the HTTP response code used is 303. See https://en.wikipedia.org/wiki/HTTP_303 for more information.
	} else {
		c.HTML(200, "login.tmpl", gin.H{"loggedIn": loggedIn, "HeaderData": HeaderData})
	}
}
