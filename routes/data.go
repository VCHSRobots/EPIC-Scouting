package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/calc"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wcharczuk/go-chart"
)

//Data route for data display
func Data(c *gin.Context) {
	querydisplay := c.Query("display")
	team, _ := strconv.Atoi(c.Query("team"))
	HeaderData := &web.HeaderData{Title: "Data", StyleSheets: []string{"global"}}
	if querydisplay == "match" {
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "MatchData": true})
	} else if querydisplay == "team" {
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "TeamOverall": true})
	} else if querydisplay == "teamprofile" {
		testTeamID := "0b28675e-4dbd-413b-96ca-016be82c78d6"
		campaign, _ := db.GetTeamCampaign(testTeamID)
		overall := calc.TeamOverall(team, campaign)
		auto := calc.TeamAuto(team, campaign)
		shooting := calc.TeamAuto(team, campaign)
		colorwheel := calc.TeamColorWheel(team, campaign)
		climbing := calc.TeamClimbing(team, campaign)
		fouls := calc.TeamFoul(team, campaign)
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "TeamProfile": true, "Overall": overall, "Auto": auto, "Shooting": shooting, "ColorWheel": colorwheel, "Climbing": climbing, "Fouls": fouls})
	} else {
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "none": true})
	}
}

//MatchDataGet sends match data in csv form to the ajax frontend
func MatchDataGet(c *gin.Context) {
	teamSortKeys := []string{"team", "overall", "auto", "shooting", "climing", "colorwheel", "fouls"}
	sortby := c.Query("sortby")
	userTeam, _ := strconv.Atoi(auth.CheckTeam(c))
	userTeamID, _ := db.GetTeamID(userTeam)
	campaign, _ := db.GetTeamCampaign(userTeamID)
	if sortby == "" || !contains(teamSortKeys, sortby) {
		sortby = "match"
	}
	//TODO: get each match from database and sort based on the querystring
	data, _ := db.GetTeamResults(userTeam, campaign)
	fmt.Println(data)
}

func contains(arr []string, val string) bool {
	for _, x := range arr {
		if x == val {
			return true
		}
	}
	return false
}

//GetGraph gets a test graph
func GetGraph(c *gin.Context) {
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0},
			},
		},
	}
	imgdata := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, imgdata)
	bytes := imgdata.Bytes()
	c.Data(http.StatusOK, "bytes", bytes)
}
