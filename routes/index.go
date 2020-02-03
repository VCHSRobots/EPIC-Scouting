package routes

import (
	"EPIC-Scouting/lib/config"
	"EPIC-Scouting/lib/web"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
indexData defines the variables that may be passed to the index page.
*/
type indexData struct {
	HeaderData *web.HeaderData
	UserMode   string
	Build      string
}

/*
Index shows the index.
*/
func Index(c *gin.Context) {
	buildName, buildDate := config.BuildInformation()
	build := fmt.Sprintf("%s [%s]", buildName, buildDate)
	// TODO: Get user status via authentication function.
	headerData := &web.HeaderData{Title: "Index", StyleSheets: nil}
	data := &indexData{headerData, "sysAdmin", build}
	c.HTML(200, "index.tmpl", data)
	// TODO: If the user is already logged in, show the "DASHBOARD" and "LOG OUT" buttons and remove the "LOGIN" and "REGISTER" buttons.
}
