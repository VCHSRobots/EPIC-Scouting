package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"fmt"

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
	//*uncomment this line to create test match*
	//db.CreateMatch("epicevent", userID, 1, true)
	//team, _ := c.Cookie("team")
	//original team id do not steal
	testTeamID := "4415epicrobotz"
	if userID != "" {
		db.StoreMatch(data.Data[0], userID, testTeamID)
		data, _ := db.GetTeamScoutData(testTeamID)
		for _, match := range *data {
			fmt.Println(match.MatchID, match.MatchNum, match.Team, match.AutoLineCross, match.AutoLowBalls, match.AutoHighBalls, match.AutoBackBalls, match.AutoPickups, match.ShotQuantity, match.LowFuel, match.HighFuel, match.BackFuel, match.StageOneComplete, match.StageOneTime, match.StageTwoComplete, match.StageTwoTime, match.Fouls, match.TechFouls, match.Card, match.ClimbTime, match.Comments)
		}
	} else {
		Forbidden(c)
	}
}

//PitPOST processes and stores/updates pit scouting data for a team
