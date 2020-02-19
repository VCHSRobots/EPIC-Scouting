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
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "TeamData": true})
	} else {
		c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData, "none": true})
	}
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
