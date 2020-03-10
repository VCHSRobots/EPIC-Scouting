package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/calc"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wcharczuk/go-chart"
)

//Images struct for exporting images
type Images struct {
	Images []string `json:"images"`
}

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
		var build strings.Builder
		var comments string
		teamID, _ := db.GetTeamID(team)
		campaign, _ := db.GetTeamCampaign(teamID)
		overall := calc.TeamOverall(team, campaign)
		auto := calc.TeamAuto(team, campaign)
		shooting := calc.TeamShooting(team, campaign)
		colorwheel := calc.TeamColorWheel(team, campaign)
		climbing := calc.TeamClimbing(team, campaign)
		fouls := calc.TeamFoul(team, campaign)
		commentList, _ := db.GetTeamComments(team, campaign)
		for ind, comment := range commentList {
			build.WriteString(comment)
			if ind != len(commentList)-1 {
				build.WriteString(", ")
			}
		}
		comments = build.String()
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "TeamProfile": true, "Overall": overall, "Auto": auto, "Shooting": shooting, "ColorWheel": colorwheel, "Climbing": climbing, "Fouls": fouls, "Comments": comments})
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

/*
MatchDataGet gets list of match data to display on match data page
*/
func MatchDataGet(c *gin.Context) {
	var build strings.Builder
	var csvList []string
	var csvString string
	var matchResult calc.MatchResults
	matchResults := make([]calc.MatchResults, 0)
	userTeam, _ := strconv.Atoi(auth.CheckTeam(c))
	userTeamID, _ := db.GetTeamID(userTeam)
	campaign, _ := db.GetTeamCampaign(userTeamID)
	matchIDs := db.ListMatchIDs(campaign)
	for _, matchID := range matchIDs {
		matchResult, _ = calc.GetMatchData(matchID)
		matchResults = append(matchResults, matchResult)
	}
	for ind, result := range matchResults {
		csvList = []string{strconv.Itoa(result.MatchNum), fmt.Sprint(result.RedParticipants), fmt.Sprint(result.BlueParticipants), strconv.Itoa(result.RedAutoBalls + result.RedTeleopShots), strconv.Itoa(result.BlueAutoBalls + result.BlueTeleopShots), strconv.Itoa(result.RedShieldStage), strconv.Itoa(result.BlueShieldStage), fmt.Sprint(result.RedClimbStatus), fmt.Sprint(result.BlueClimbStatus), fmt.Sprint(result.RedRankingPoints), fmt.Sprint(result.BlueRankingPoints), strconv.Itoa(result.RedPoints), strconv.Itoa(result.BluePoints), result.Winner}
		build.WriteString(writeCSVString(csvList))
		if ind != len(matchResults)-1 {
			build.WriteString("\n")
		}
	}
	csvString = build.String()
	c.String(http.StatusOK, "text", csvString)
}

/*
TeamMatchDataGet gets match statistics for each match a team participated in
*/
func TeamMatchDataGet(c *gin.Context) {
	var build strings.Builder
	var csvList []string
	var csvString, balanced, teammates, opponents string
	var matchResult db.MatchData
	var matches []db.MatchData
	var participants [][]int
	userTeam, _ := strconv.Atoi(auth.CheckTeam(c))
	userTeamID, _ := db.GetTeamID(userTeam)
	campaign, _ := db.GetTeamCampaign(userTeamID)
	matchIDs := db.ListMatchIDs(campaign)
	teamNumString := c.Query("team")
	teamNum, _ := strconv.Atoi(teamNumString)
	for ind, matchID := range matchIDs {
		matchResult = calc.ResolveMatchConflicts(teamNum, matchID)
		matches = []db.MatchData{matchResult}
		if matchResult.Balanced {
			balanced = "true"
		} else {
			balanced = "false"
		}
		participants = db.GetMatchParticipants(matchID)
		if containsInt(participants[0], teamNum) {
			teammates = fmt.Sprint(participants[0])
			opponents = fmt.Sprint(participants[1])
		} else {
			teammates = fmt.Sprint(participants[1])
			opponents = fmt.Sprint(participants[0])
		}
		csvList = []string{strconv.Itoa(matchResult.MatchNum), strconv.Itoa(calc.Overall(matches)), teammates, opponents, strconv.Itoa(calc.Shooting(matches)), strconv.Itoa(calc.Auto(matches)), strconv.Itoa(calc.ColorWheel(matches)), matchResult.Climbed, balanced, strconv.Itoa(calc.Foul(matches))}
		build.WriteString(writeCSVString(csvList))
		if ind != len(matchIDs)-1 {
			build.WriteString("\n")
		}
	}
	csvString = build.String()
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

func containsInt(arr []int, val int) bool {
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

func writeCSVString(arr []string) string {
	var str strings.Builder
	for ind, val := range arr {
		str.WriteString(val)
		if ind != len(arr)-1 {
			str.WriteString(",")
		}
	}
	return str.String()
}

/*
GetMatchHistory gets table with a team's match history
*/
func GetMatchHistory(c *gin.Context) {

}

/*
GetTeamImages gets a team's images
*/
func GetTeamImages(c *gin.Context) {
	var images Images
	teamID, _ := db.GetTeamID(4415)
	campaignID, _ := db.GetTeamCampaign(teamID)
	teamNum, _ := strconv.Atoi(c.Query("team"))
	imageList, _ := db.GetTeamImages(teamNum, campaignID)
	images.Images = imageList
	jsonBytes, _ := json.Marshal(images)
	c.Data(http.StatusOK, "json", jsonBytes)
}

//GetGraph gets a graph which responds to querystring parameters
func GetGraph(c *gin.Context) {
	var xAxis, yAxis string
	matchGroups := make(map[int][]db.MatchData, 0)
	x := make([]float64, 1)
	y := make([]float64, 1)
	graphSubject := c.Query("subject")
	userTeam := c.Query("team")
	userTeamNum, _ := strconv.Atoi(userTeam)
	userTeamID, _ := db.GetTeamID(userTeamNum)
	campaign, _ := db.GetTeamCampaign(userTeamID)
	if graphSubject == "Overall" {
		xAxis = c.Query("team")
		yAxis = "Overall"
		teamNum, _ := strconv.Atoi(xAxis)
		matches, _ := db.GetTeamMatches(teamNum, campaign)
		for _, match := range *matches {
			_, ok := matchGroups[match.MatchNum]
			if ok {
				matchGroups[match.MatchNum] = append(matchGroups[match.MatchNum], match)
			} else {
				matchGroups[match.MatchNum] = make([]db.MatchData, 1)
				matchGroups[match.MatchNum][0] = match
			}
		}
		for matchNum, matches := range matchGroups {
			x = append(x, float64(matchNum))
			y = append(y, float64(calc.Overall(matches)))
		}
	} else if graphSubject == "Auto" {
		xAxis = c.Query("team")
		yAxis = "Auto"
		teamNum, _ := strconv.Atoi(xAxis)
		matches, _ := db.GetTeamMatches(teamNum, campaign)
		for _, match := range *matches {
			_, ok := matchGroups[match.MatchNum]
			if ok {
				matchGroups[match.MatchNum] = append(matchGroups[match.MatchNum], match)
			} else {
				matchGroups[match.MatchNum] = make([]db.MatchData, 1)
				matchGroups[match.MatchNum][0] = match
			}
		}
		for matchNum, matches := range matchGroups {
			x = append(x, float64(matchNum))
			y = append(y, float64(calc.Auto(matches)))
		}
	} else if graphSubject == "Shooting" {
		xAxis = c.Query("team")
		yAxis = "Shooting"
		teamNum, _ := strconv.Atoi(xAxis)
		matches, _ := db.GetTeamMatches(teamNum, campaign)
		for _, match := range *matches {
			_, ok := matchGroups[match.MatchNum]
			if ok {
				matchGroups[match.MatchNum] = append(matchGroups[match.MatchNum], match)
			} else {
				matchGroups[match.MatchNum] = make([]db.MatchData, 1)
				matchGroups[match.MatchNum][0] = match
			}
		}
		for matchNum, matches := range matchGroups {
			x = append(x, float64(matchNum))
			y = append(y, float64(calc.Shooting(matches)))
		}
	} else if graphSubject == "ColorWheel" {
		xAxis = c.Query("team")
		yAxis = "Color Wheel"
		teamNum, _ := strconv.Atoi(xAxis)
		matches, _ := db.GetTeamMatches(teamNum, campaign)
		for _, match := range *matches {
			_, ok := matchGroups[match.MatchNum]
			if ok {
				matchGroups[match.MatchNum] = append(matchGroups[match.MatchNum], match)
			} else {
				matchGroups[match.MatchNum] = make([]db.MatchData, 1)
				matchGroups[match.MatchNum][0] = match
			}
		}
		for matchNum, matches := range matchGroups {
			x = append(x, float64(matchNum))
			y = append(y, float64(calc.ColorWheel(matches)))
		}
	} else if graphSubject == "Climbing" {
		xAxis = c.Query("team")
		yAxis = "Climbing"
		teamNum, _ := strconv.Atoi(xAxis)
		matches, _ := db.GetTeamMatches(teamNum, campaign)
		for _, match := range *matches {
			_, ok := matchGroups[match.MatchNum]
			if ok {
				matchGroups[match.MatchNum] = append(matchGroups[match.MatchNum], match)
			} else {
				matchGroups[match.MatchNum] = make([]db.MatchData, 1)
				matchGroups[match.MatchNum][0] = match
			}
		}
		for matchNum, matches := range matchGroups {
			x = append(x, float64(matchNum))
			y = append(y, float64(calc.Climbing(matches)))
		}
	} else if graphSubject == "Fouls" {
		xAxis = c.Query("team")
		yAxis = "Fouls"
		teamNum, _ := strconv.Atoi(xAxis)
		matches, _ := db.GetTeamMatches(teamNum, campaign)
		for _, match := range *matches {
			_, ok := matchGroups[match.MatchNum]
			if ok {
				matchGroups[match.MatchNum] = append(matchGroups[match.MatchNum], match)
			} else {
				matchGroups[match.MatchNum] = make([]db.MatchData, 1)
				matchGroups[match.MatchNum][0] = match
			}
		}
		for matchNum, matches := range matchGroups {
			x = append(x, float64(matchNum))
			y = append(y, float64(calc.Foul(matches)))
		}
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: xAxis,
		},
		YAxis: chart.YAxis{
			Name: yAxis,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: x,
				YValues: y,
			},
		},
	}
	imgdata := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, imgdata)
	bytes := imgdata.Bytes()
	c.Data(http.StatusOK, "bytes", bytes)
}
