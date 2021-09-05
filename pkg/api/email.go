package api

import (
	"github.com/gin-gonic/gin"
	"log"
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
		log.Println("Failed to get user active code")
		return
	}

	err := srv.RunUserActive(code)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		log.Println(err.Error())
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

func emailChange(c *gin.Context) {
	errMsg := "Failed to change email"
	code, ok := c.GetQuery("verify-code")
	if !ok {
		c.String(http.StatusBadRequest, errMsg)
		log.Println("Failed to get email edit code")
		return
	}

	err := srv.RunEmailEdit(code)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		log.Println(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

type emailRecoverReq struct {
	Password string `json:"password" validate:"omitempty,max=255,min=8"`
}

func emailRecover(c *gin.Context) {
	req := emailRecoverReq{}
	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	errMsg := "Failed to reset the password"
	code, ok := c.GetQuery("recover-code")
	if !ok {
		c.String(http.StatusBadRequest, errMsg)
		log.Println("Failed to get user recover code")
		return
	}

	pwd, err := srv.GenPwd(req.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	err = srv.RunUserRecover(code, pwd)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		log.Println(err.Error())
		return
	}

	c.Status(http.StatusOK)
}
