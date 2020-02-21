package routes

import (
	"EPIC-Scouting/lib/web"
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wcharczuk/go-chart"
)

//Data route for data display
func Data(c *gin.Context) {
	querydisplay := c.Query("display")
	HeaderData := &web.HeaderData{Title: "Data", StyleSheets: []string{"global"}}
	if querydisplay == "match" {
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "MatchData": true})
	} else if querydisplay == "team" {
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "TeamOverall": true})
	} else if querydisplay == "teamprofile" {
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "TeamProfile": true})
	} else {
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "none": true})
	}
}

//MatchDataGet sends match data in csv form to the ajax frontend
func MatchDataGet(c *gin.Context) {
	teamSortKeys := []string{"team", "overall", "auto", "shooting", "climing", "colorwheel", "fouls"}
	sortby := c.Query("sortby")
	if sortby == "" || !contains(teamSortKeys, sortby) {
		sortby = "match"
	}
	//TODO: get each match from database and sort based on the querystring
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
