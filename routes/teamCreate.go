package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
TeamCreate shows the team creation page.
*/
func TeamCreate(c *gin.Context) {
	HeaderData := &web.HeaderData{Title: "Create Team", StyleSheets: []string{"global"}}
	c.HTML(http.StatusOK, "teamCreate.tmpl", gin.H{"HeaderData": HeaderData})
}

/*
TeamCreatePOST TODO
*/
func TeamCreatePOST(c *gin.Context) {
	c.Request.ParseForm()
	teamCreator := auth.CheckLogin(c)
	if teamCreator == "" {
		Forbidden(c)
	}
	teamNum, _ := strconv.Atoi(c.PostForm("number"))
	teamName := c.PostForm("name")
	db.TeamCreate(teamNum, teamName, "")
}
