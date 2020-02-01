package routes

import (
	"github.com/gin-gonic/gin"
)

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
	/*
		// TODO: Make this actually work!
		var ctx Credentials
		err := c.Bind(&ctx)
		if err != nil {
			print("RIP")
		}
	*/
	c.Request.ParseForm()
	username := c.PostForm("username")
	print("USR: " + string(username))
	password := c.PostForm("password")
	print("PWD: " + string(password))
}
