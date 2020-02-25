package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
TeamJoin shows the join team page.
*/
func TeamJoin(c *gin.Context) {
	c.HTML(http.StatusOK, "teamJoin.tmpl", nil)
}

/*
TeamJoinRequest requests to join a team.
*/
func TeamJoinRequest(c *gin.Context) {
	//db.TeamJoinRequest(requesterID, teamID string)
}
