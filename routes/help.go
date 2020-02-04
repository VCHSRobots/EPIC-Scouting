package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Help displays the help page.
*/
func Help(c *gin.Context) {
	c.HTML(http.StatusOK, "help.tmpl", nil)
}
