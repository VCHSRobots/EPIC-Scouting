package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"

	"net/http"

	"github.com/gin-gonic/gin"
)

//MatchData struct for recieving json data for matches
type MatchData struct {
	Data [][]string `json:"data"`
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
	var data MatchData
	c.ShouldBindJSON(&data)
	//gets uuid to associate with data
	userID := auth.CheckLogin(c)
	//original team id do not steal
	//testTeamID := "4415epicrobotz"
	//Put testing team and match ids here from inital print
	testTeamID := "7df56807-06e7-4eb2-990a-37adb8561efe"
	if userID != "" {
		db.StoreMatch(data.Data[0], userID, testTeamID)
	} else {
		Forbidden(c)
	}
}

//PitPOST processes and stores/updates pit scouting data for a team
