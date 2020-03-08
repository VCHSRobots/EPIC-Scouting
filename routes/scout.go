package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"

	"net/http"

	"github.com/gin-gonic/gin"
)

//PostData struct for receiving json data for all post functions
type PostData struct {
	Data []string `json:"data"`
}

/*
Scout shows the scout page.
*/
func Scout(c *gin.Context) {
	querytype := c.Query("type")
	if querytype == "match" {
		HeaderData := &web.HeaderData{Title: "Match Scouting", StyleSheets: []string{"scout"}}
		c.HTML(http.StatusOK, "scout.tmpl", gin.H{"HeaderData": HeaderData, "MatchScout": true})
	} else if querytype == "pit" {
		HeaderData := &web.HeaderData{Title: "Pit Scouting", StyleSheets: []string{"scout"}}
		c.HTML(http.StatusOK, "scout.tmpl", gin.H{"HeaderData": HeaderData, "PitScout": true})
	} else {
		HeaderData := &web.HeaderData{Title: "Scouting?", StyleSheets: []string{"scout"}}
		c.HTML(http.StatusOK, "scout.tmpl", gin.H{"HeaderData": HeaderData, "nope": true})
	}
}

//MatchPOST processes and stores scouting data from a match
func MatchPOST(c *gin.Context) {
	var data PostData
	c.ShouldBindJSON(&data)
	//gets uuid to associate with data
	userID := auth.CheckLogin(c)
	//TODO: automatically gets 4415's team id - change in future
	teamID, _ := db.GetTeamID(4415)
	if userID != "" {
		db.StoreMatch(data.Data, userID, teamID)
	} else {
		Forbidden(c)
	}
}

//PitPOST processes and stores/updates pit scouting data for a team
func PitPOST(c *gin.Context) {
	var data PostData
	c.ShouldBindJSON(&data)
	userID := auth.CheckLogin(c)
	teamID, _ := db.GetTeamID(4415)
	if userID != "" {
		db.WritePitData(data.Data, userID, teamID)
	} else {
		Forbidden(c)
	}
}
