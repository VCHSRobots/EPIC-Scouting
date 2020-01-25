package routes

import (
	"github.com/gin-gonic/gin"
)

/*
Scout shows the scout page.
*/
func Scout(c *gin.Context) {
	c.HTML(200, "scout.tmpl", nil)
}
