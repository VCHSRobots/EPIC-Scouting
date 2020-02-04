package routes

import (
	"EPIC-Scouting/lib/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Login shows the login page.
*/
func Login(c *gin.Context) {
	loggedIn := true
	c.HTML(200, "login.tmpl", loggedIn)
}

/*
LoginPOST logs a user in.
*/
func LoginPOST(c *gin.Context) {
	c.Request.ParseForm()
	username := c.PostForm("username")
	password := c.PostForm("password")
	loggedIn, userdata := db.CheckLogin(username, password)
	fmt.Println(loggedIn, userdata)
	if loggedIn {
		c.Redirect(http.StatusSeeOther, "/dashboard") // Although gin's method here is named Redirect, the HTTP response code used is 303. See https://en.wikipedia.org/wiki/HTTP_303 for more information.
	} else {
		c.HTML(200, "login.tmpl", loggedIn)
	}
}
