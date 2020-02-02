package routes

import (
	"EPIC-Scouting/lib/db"
	"fmt"

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
	println("USR: " + string(username))
	password := c.PostForm("password")
	println("PWD: " + string(password))
	worked, userdata := db.CheckLogin(username, password)
	fmt.Println(worked, userdata)
	if worked {
		c.HTML(200, "dashboard.tmpl", nil)
	} else {
		c.HTML(200, "login.tmpl", nil)
	}
}
