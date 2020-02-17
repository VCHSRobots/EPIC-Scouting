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
	HeaderData := &web.HeaderData{Title: "Data", StyleSheets: []string{"global"}}
	c.HTML(http.StatusOK, "data.tmpl", gin.H{"HeaderData": HeaderData})
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
