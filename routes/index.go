package routes

import (
	"EPIC-Scouting/lib/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
IndexData defines the variables that may be passed to the index page.
*/
type IndexData struct {
	Title    string
	UserMode string
	Build    string
}

/*
Index shows the index.
*/
func Index(c *gin.Context) {
	buildName, buildDate := config.BuildInformation()
	build := fmt.Sprintf("%s [%s]", buildName, buildDate)
	// TODO: Get user status via authentication function.
	data := &IndexData{"le title", "guest", build} // TODO: Header needs to be re-executed in order to properly use Title. This should be on a per-page basis.
	c.HTML(200, "index.tmpl", data)
	// TODO: If the user is already logged in, show the "DASHBOARD" and "LOG OUT" buttons and remove the "LOGIN" and "REGISTER" buttons.
}
