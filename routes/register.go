package routes

import (
	"EPIC-Scouting/lib/db"
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
	c.Request.ParseForm()
	var d db.UserData
	d.UserName = c.PostForm("username")
	d.Password = c.PostForm("password")
	d.Email = c.PostForm("email")
	d.FirstName = c.PostForm("firstname")
	d.LastName = c.PostForm("lastname")
	created, err := db.UserCreate(&d)
	if created {
		c.HTML(200, "registered.tmpl", nil)
	} else {
		c.HTML(200, "register.tmpl", gin.H{"Error": err.Error()})
	}
}
