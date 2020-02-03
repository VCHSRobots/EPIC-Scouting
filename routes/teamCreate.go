package routes

import (
	"github.com/gin-gonic/gin"
)

/*
TeamCreate shows the team creation page.
*/
func TeamCreate(c *gin.Context) {
	c.HTML(200, "teamCreate.tmpl", nil)
}
