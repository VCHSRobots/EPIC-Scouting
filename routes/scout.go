package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Scout shows the scout page.
*/
func Scout(c *gin.Context) {
	c.HTML(http.StatusOK, "scout.tmpl", nil)
}

//MatchScoutPOST processes and stores scouting data from a match

//PitScoutPOST processes and stores/updates pit scoutin data for a team
