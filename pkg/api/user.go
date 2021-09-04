package api

import (
	"github.com/gin-gonic/gin"
	"mlauth/pkg/dao"
	"mlauth/pkg/mdl"
	"mlauth/pkg/srv"
	"net/http"
	"time"
)

type userDetail struct {
	Uid         int       `json:"uid"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

func userGet(c *gin.Context) {
	u := c.MustGet("user").(mdl.User)
	res := userDetail{
		Uid:         u.Uid,
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt,
	}
	c.JSON(http.StatusOK, res)
}

type userRegisterReq struct {
	Username    string `json:"username" validate:"max=255,min=2,alphanum"`
	Password    string `json:"password" validate:"max=255,min=8"`
	Email       string `json:"email" validate:"email"`
	DisplayName string `json:"display_name" validate:"max=255"`
}

func userRegister(c *gin.Context) {
	req := userRegisterReq{}
	err := c.Bind(&req)
	if err != nil {
		return
	}

	pwd, err := srv.GenPwd(req.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	uCreate := mdl.User{
		Username:    req.Username,
		Password:    pwd,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		IsActive:    false,
		IsSuper:     false,
		CreatedAt:   time.Now().UTC(),
	}
	u, err := dao.InsertUser(uCreate)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if !dao.CheckEmailRetry("user-active", u.Uid) {
		c.String(http.StatusBadRequest, "Email request too often")
		return
	}

	go func() {
		_ = srv.ReqUserActive(u)
	}()

	c.JSON(http.StatusOK, userMdl2userDetail(u))
}

type userEditReq struct {
	DisplayName string `json:"display_name" validate:"omitempty,max=255"`
	Password    string `json:"password" validate:"omitempty,max=255,min=8"`
	Email       string `json:"email" validate:"omitempty,email"`
}

func userEdit(c *gin.Context) {
	uPre := c.MustGet("user").(mdl.User)
	req := userEditReq{}
	err := c.Bind(&req)
	if err != nil {
		return
	}

	uEdit := uPre
	if req.DisplayName != "" {
		uEdit.DisplayName = req.DisplayName
	}
	if req.Password != "" {
		pwd, err := srv.GenPwd(req.Password)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		uEdit.Password = pwd
	}
	if req.Email != "" {
		if !dao.CheckEmailRetry("email-edit", uPre.Uid) {
			c.String(http.StatusBadRequest, "Email request too often")
			return
		}

		go func() {
			_ = srv.ReqEmailEdit(uPre, req.Email)
		}()
	}

	u, err := dao.UpdateUser(uPre.Uid, uEdit)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, userMdl2userDetail(u))
}

func userMdl2userDetail(u mdl.User) userDetail {
	return userDetail{
		Uid:         u.Uid,
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt,
	}
}
