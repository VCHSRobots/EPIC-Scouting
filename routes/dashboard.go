package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Dashboard shows the dashboard.
*/
func Dashboard(c *gin.Context) {
	userID := auth.CheckLogin(c)
	println(userID)
	if userID == "" {
		Forbidden(c)
		return
	}
	userData, _ := db.UserQuery(userID)
	HeaderData := &web.HeaderData{Title: "Dashboard", StyleSheets: []string{"global"}}
	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{"HeaderData": HeaderData, "uuid": userID, "Username": userData.UserName, "SysAdmin": userData.SysAdmin})
}
