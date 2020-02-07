package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
TeamAdmin shows the team administration page.
*/
func TeamAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "teamAdmin.tmpl", nil)
}
