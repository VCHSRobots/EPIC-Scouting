package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
TeamData shows the team data page.
*/
func TeamData(c *gin.Context) {
	c.HTML(http.StatusOK, "teamData.tmpl", nil)
}
