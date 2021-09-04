package api

import (
	"github.com/gin-gonic/gin"
	"mlauth/pkg/dao"
	"mlauth/pkg/mdl"
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

func emailActiveRetry(c *gin.Context) {
	u := c.MustGet("user").(mdl.User)

	if !dao.CheckEmailRetry("user-active", u.Uid) {
		c.String(http.StatusBadRequest, "Email request too often")
		return
	}

	go func() {
		_ = srv.ReqUserActive(u)
	}()

	c.Status(http.StatusOK)
}
