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
