package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/calc"
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
	//testTeamID := "4415epicrobotz"
	//Put testing team and match ids here from inital print
	testTeamID := "0b28675e-4dbd-413b-96ca-016be82c78d6"
	testMatchID := "957d87fd-d3de-46d8-8c5b-ed4408ca738b"
	campaign, _ := db.GetTeamCampaign(testTeamID)
	if userID != "" {
		db.StoreMatch(data.Data[0], userID, testTeamID)
		data, _ := db.GetTeamResults(4415, testMatchID)
		if data != nil {
			for _, match := range *data {
				fmt.Println(match.MatchID, match.MatchNum, match.Team, match.AutoLineCross, match.AutoLowBalls, match.AutoHighBalls, match.AutoBackBalls, match.AutoPickups, match.ShotQuantity, match.LowFuel, match.HighFuel, match.BackFuel, match.StageOneComplete, match.StageOneTime, match.StageTwoComplete, match.StageTwoTime, match.Fouls, match.TechFouls, match.Card, match.ClimbTime, match.Comments)
			}
			calc.TeamOverall(4415, campaign)
		}
	} else {
		Forbidden(c)
	}
}

//PitPOST processes and stores/updates pit scouting data for a team
