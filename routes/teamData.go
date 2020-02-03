package routes

import (
	"github.com/gin-gonic/gin"
)

/*
TeamData shows the team data page.
*/
func TeamData(c *gin.Context) {
	c.HTML(200, "teamData.tmpl", nil)
}
