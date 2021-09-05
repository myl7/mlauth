package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"mlauth/pkg/dao"
	"mlauth/pkg/mdl"
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
	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	errMsg := "Failed to login"
	u, err := dao.SelectUserByUsername(req.Username)
	if err != nil {
		c.String(http.StatusForbidden, errMsg)
		log.Println(err.Error())
		return
	}

	if !srv.CheckPwd(u.Password, req.Password) {
		c.String(http.StatusForbidden, errMsg)
		log.Println("Password mismatch")
		return
	}

	at, err := srv.GenAccessToken(u.Uid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	ut, err := srv.GenUpdateToken(u.Uid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err.Error())
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
	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	errMsg := "Failed to renew access token"
	uid, err := srv.CheckUpdateToken(req.UpdateToken)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		log.Println(err.Error())
		return
	}

	u, err := dao.SelectUser(uid)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		log.Println(err.Error())
		return
	}

	at, err := srv.GenAccessToken(u.Uid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	res := userRenewRes{
		UpdateToken: req.UpdateToken,
		AccessToken: at,
	}
	c.JSON(http.StatusOK, res)
}

type userRecoverReq struct {
	Username string `json:"username" validate:"max=255,min=2,alphanum"`
}

func userRecover(c *gin.Context) {
	req := userRecoverReq{}
	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	u, err := dao.SelectUserByUsername(req.Username)
	if err != nil {
		c.Status(http.StatusOK)
		log.Println(err.Error())
		return
	}

	err = dao.SetEmailRetry("user-recover", u.Uid)
	if err != nil {
		c.String(http.StatusBadRequest, "Email request too often")
		return
	}

	go func() {
		_ = srv.ReqUserRecover(u)
	}()

	c.Status(http.StatusOK)
}

func userAuth(c *gin.Context) {
	u, err := getUserInAuth(c)
	if err != nil {
		return
	}

	if !u.IsActive {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Set("user", u)
}

func userAuthExist(c *gin.Context) {
	u, err := getUserInAuth(c)
	if err != nil {
		return
	}

	c.Set("user", u)
}

func getUserInAuth(c *gin.Context) (mdl.User, error) {
	at := c.GetHeader("x-access-token")
	uid, err := srv.CheckAccessToken(at)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return mdl.User{}, err
	}

	u, err := dao.SelectUser(uid)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return mdl.User{}, err
	}

	return u, nil
}
