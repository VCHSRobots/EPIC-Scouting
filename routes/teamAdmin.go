package routes

import (
	"github.com/gin-gonic/gin"
)

/*
TeamAdmin shows the team administration page.
*/
func TeamAdmin(c *gin.Context) {
	c.HTML(200, "teamAdmin.tmpl", nil)
}
