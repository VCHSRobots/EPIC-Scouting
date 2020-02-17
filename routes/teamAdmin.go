package routes

import (
	"EPIC-Scouting/lib/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
TeamAdmin shows the team administration page.
*/
func TeamAdmin(c *gin.Context) {
	HeaderData := &web.HeaderData{Title: "Team Admin", StyleSheets: []string{"global"}}
	c.HTML(http.StatusOK, "teamAdmin.tmpl", gin.H{"HeaderData": HeaderData})
}
