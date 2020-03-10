package routes

import (
	"EPIC-Scouting/lib/config"
	"EPIC-Scouting/lib/web"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
indexData defines the variables that may be passed to the index page.
*/
type indexData struct {
	HeaderData *web.HeaderData
	UserMode   string
	Build      string
}

/*
Index shows the index.
*/
func Index(c *gin.Context) {
	buildName, buildDate := config.BuildInformation()
	build := fmt.Sprintf("%s [%s]", buildName, buildDate)
	// TODO: Get user status via authentication function.
	headerData := &web.HeaderData{Title: "Index", StyleSheets: []string{"index"}}
	data := &indexData{headerData, "guest", build}
	c.HTML(http.StatusOK, "index.tmpl", data)
}
