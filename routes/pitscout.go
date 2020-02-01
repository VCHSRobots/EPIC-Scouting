package routes

import (
	"github.com/gin-gonic/gin"
)

/*
PitScout shows the scout page.
*/
func PitScout(c *gin.Context) {
	c.HTML(200, "pitscout.tmpl", nil)
}
