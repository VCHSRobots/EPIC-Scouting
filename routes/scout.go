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
