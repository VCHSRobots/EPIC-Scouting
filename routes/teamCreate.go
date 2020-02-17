package routes

import (
	"EPIC-Scouting/lib/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
TeamCreate shows the team creation page.
*/
func TeamCreate(c *gin.Context) {
	HeaderData := &web.HeaderData{Title: "Create Team", StyleSheets: []string{"global"}}
	c.HTML(http.StatusOK, "teamCreate.tmpl", gin.H{"HeaderData": HeaderData})
}
