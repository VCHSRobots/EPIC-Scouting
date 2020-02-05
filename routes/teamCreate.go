package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
TeamCreate shows the team creation page.
*/
func TeamCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "teamCreate.tmpl", nil)
}
