package routes

import "github.com/gin-gonic/gin"

/*
TeamJoin shows the join team page.
*/
func TeamJoin(c *gin.Context) {
	c.HTML(200, "teamJoin.tmpl", nil)
}
