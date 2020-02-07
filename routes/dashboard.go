package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Dashboard shows the dashboard.
*/
func Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.tmpl", nil)
}
