package routes

import (
	"EPIC-Scouting/lib/db"
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

/*
TeamCreatePOST TODO
*/

func TeamCreatePOST(c *gin.Context) {
	c.Request.ParseForm()
	teamCreator := "USERID-GOES-HERE" // TODO: Get UserID via session cookie.
	c.PostForm("number")
	c.PostForm("name")
	db.TeamCreate()
}
