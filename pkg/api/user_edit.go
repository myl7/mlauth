package api

import (
	"github.com/gin-gonic/gin"
	"mlauth/pkg/dao"
	"mlauth/pkg/mdl"
	"mlauth/pkg/srv"
	"net/http"
)

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
		// TODO
	}

	u, err := dao.UpdateUser(uPre.Uid, uEdit)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, userMdl2userDetail(u))
}
