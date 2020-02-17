package routes

import (
	"EPIC-Scouting/lib/web"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
Dashboard shows the dashboard.
*/
func Dashboard(c *gin.Context) {
	cookiestring, _ := c.Cookie("login")
	if cookiestring != "" {
		cstrs := strings.Fields(cookiestring)
		HeaderData := &web.HeaderData{Title: "Dashboard", StyleSheets: []string{"global"}}
		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{"HeaderData": HeaderData, "uuid": cstrs[0], "username": cstrs[1]})
	} else {
		Forbidden(c)
	}

}
