package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Register shows the register page.
*/
func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", nil)
}

/*
RegisterPOST processes user registration form
*/
func RegisterPOST(c *gin.Context) {
	/*c.Request.ParseForm()
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	firstName := c.PostForm("firstname")
	lastName := c.PostForm("lastname")
	phone := c.PostForm("phone")
	db.CreateUser(db.DatabasePath, username, password, firstName, lastName, email, phone, "user")
	log.Debugf("Created user: ")
	c.HTML(200, "registered.tmpl", nil)*/
}
