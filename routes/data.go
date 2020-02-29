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
	"strings"

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
		teamID, _ := db.GetTeamID(team)
		campaign, _ := db.GetTeamCampaign(teamID)
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

//TeamDataGet sends match data in csv form to the ajax frontend
func TeamDataGet(c *gin.Context) {
	var swap []int
	var build strings.Builder
	teamSortKeys := []string{"Team", "Overall", "Auto", "Shooting", "Climing", "Colorwheel", "Fouls"}
	sortby := c.Query("sortby")
	userTeam, _ := strconv.Atoi(auth.CheckTeam(c))
	userTeamID, _ := db.GetTeamID(userTeam)
	campaign, _ := db.GetTeamCampaign(userTeamID)
	if sortby == "" || !contains(teamSortKeys, sortby) {
		sortby = "Overall"
	}
	searchind := where(teamSortKeys, sortby)
	fmt.Println(sortby, searchind)
	//TODO: get each match from database and sort based on the querystring
	scores := calc.GetTeamScores(campaign)
	for x := len(scores) - 1; x >= 0; x-- {
		for y := x - 1; y >= 0; y-- {
			if scores[y][searchind] < scores[x][searchind] {
				swap = scores[x]
				scores[x] = scores[y]
				scores[y] = swap
			}
		}
	}
	for ind, score := range scores {
		build.WriteString(writeCSV(score))
		if ind != len(score)-1 {
			build.WriteString("\n")
		}
	}
	csvString := build.String()
	c.String(http.StatusOK, "text", csvString)
}

func contains(arr []string, val string) bool {
	for _, x := range arr {
		if x == val {
			return true
		}
	}
	return false
}

func where(arr []string, val string) int {
	for ind, x := range arr {
		if x == val {
			return ind
		}
	}
	return -1
}

func writeCSV(arr []int) string {
	var str strings.Builder
	for ind, val := range arr {
		strval := strconv.Itoa(val)
		str.WriteString(strval)
		if ind != len(arr)-1 {
			str.WriteString(",")
		}
	}
	return str.String()
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
