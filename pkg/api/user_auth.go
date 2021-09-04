package api

import (
	"github.com/gin-gonic/gin"
	"mlauth/pkg/dao"
	"mlauth/pkg/srv"
	"net/http"
)

type userLoginReq struct {
	Username string `json:"username" validate:"max=255,min=2,alphanum"`
	Password string `json:"password" validate:"max=255,min=8"`
}

type userLoginRes struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
	UpdateToken string `json:"update_token"`
}

func userLogin(c *gin.Context) {
	req := userLoginReq{}
	err := c.Bind(&req)
	if err != nil {
		return
	}

	errMsg := "Failed to login"
	u, err := dao.SelectUserByUsername(req.Username)
	if err != nil {
		c.String(http.StatusForbidden, errMsg)
		return
	}

	if !srv.CheckPwd(u.Password, req.Password) {
		c.String(http.StatusForbidden, errMsg)
		return
	}

	at, err := srv.GenAccessToken(u.Uid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	ut, err := srv.GenUpdateToken(u.Uid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	res := userLoginRes{
		Username:    req.Username,
		AccessToken: at,
		UpdateToken: ut,
	}
	c.JSON(http.StatusOK, res)
}

type userRenewReq struct {
	UpdateToken string `json:"update_token"`
}

type userRenewRes struct {
	UpdateToken string `json:"update_token"`
	AccessToken string `json:"access_token"`
}

func userRenew(c *gin.Context) {
	req := userRenewReq{}
	err := c.Bind(&req)
	if err != nil {
		return
	}

	errMsg := "Failed to renew access token"
	uid, err := srv.CheckUpdateToken(req.UpdateToken)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		return
	}

	u, err := dao.SelectUser(uid)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		return
	}

	at, err := srv.GenAccessToken(u.Uid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	res := userRenewRes{
		UpdateToken: req.UpdateToken,
		AccessToken: at,
	}
	c.JSON(http.StatusOK, res)
}

func userAuth(c *gin.Context) {
	at := c.GetHeader("x-access-token")
	uid, err := srv.CheckAccessToken(at)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	u, err := dao.SelectUser(uid)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", u)
}
