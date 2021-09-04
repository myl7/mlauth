package api

import (
	"github.com/gin-gonic/gin"
	"mlauth/pkg/srv"
	"net/http"
)

func emailActive(c *gin.Context) {
	errMsg := "Failed to activate the user"
	code, ok := c.GetQuery("active-code")
	if !ok {
		c.String(http.StatusBadRequest, errMsg)
		return
	}

	err := srv.RunUserActive(code)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		return
	}

	c.Status(http.StatusOK)
}
