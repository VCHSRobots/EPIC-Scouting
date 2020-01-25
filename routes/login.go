package routes

import (
	"github.com/gin-gonic/gin"
)

/*
Credentials is a struct defining the username and password variables for LoginPOST.
*/
type Credentials struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

/*
Login shows the login page.
*/
func Login(c *gin.Context) {
	c.HTML(200, "login.tmpl", nil)
}

/*
LoginPOST logs a user in.
*/
func LoginPOST(c *gin.Context) {
	// TODO: Make this actually work!
	var ctx Credentials
	err := c.Bind(&ctx)
	if err != nil {
		print("RIP")
	}
	username := c.PostForm("username")
	print("USR: " + string(username))
	password := c.PostForm("password")
	print("PWD: " + string(password))
}
