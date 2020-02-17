package routes

import (
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Register shows the register page.
*/
func Register(c *gin.Context) {
	HeaderData := &web.HeaderData{Title: "Register", StyleSheets: []string{"global"}}
	c.HTML(http.StatusOK, "register.tmpl", gin.H{"HeaderData": HeaderData})
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
	HeaderData := &web.HeaderData{Title: "Registered", StyleSheets: []string{"global"}}
	if created {
		c.HTML(200, "registered.tmpl", gin.H{"HeaderData": HeaderData})
	} else {
		c.HTML(200, "register.tmpl", gin.H{"HeaderData": HeaderData, "Error": err.Error()})
	}
}
