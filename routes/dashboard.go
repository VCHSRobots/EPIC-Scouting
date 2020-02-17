package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Dashboard shows the dashboard.
*/
func Dashboard(c *gin.Context) {
	userMode := auth.GetUserMode(c)
	uuid, username := auth.DecodeLoginCookie(c)
	if userMode == "sysadmin" {
		HeaderData := &web.HeaderData{Title: "Dashboard", StyleSheets: []string{"global"}}
		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{"HeaderData": HeaderData, "uuid": uuid, "Username": username, "SysAdmin": true})
	} else if userMode == "user" {
		HeaderData := &web.HeaderData{Title: "Dashboard", StyleSheets: []string{"global"}}
		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{"HeaderData": HeaderData, "uuid": uuid, "Username": username})
	} else {
		Forbidden(c)
	}
}
