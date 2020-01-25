package routes

import (
	"github.com/gin-gonic/gin"
)

/*
Dashboard shows the dashboard.
*/
func Dashboard(c *gin.Context) {
	c.HTML(200, "dashboard.tmpl", nil)
}
