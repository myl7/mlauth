package api

import (
	"github.com/gin-gonic/gin"
	"mlauth/pkg/dao"
	"mlauth/pkg/mdl"
	"net/http"
)

type userEditReq struct {
	DisplayName string `json:"display_name" validate:"max=255"`
}

func userEdit(c *gin.Context) {
	uPre := c.MustGet("user").(mdl.User)
	req := userEditReq{}
	err := c.Bind(&req)
	if err != nil {
		return
	}

	uEdit := uPre
	uEdit.DisplayName = req.DisplayName
	u, err := dao.UpdateUser(uPre.Uid, uEdit)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, userMdl2userDetail(u))
}
