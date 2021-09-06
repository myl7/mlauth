package rpc

import (
	"context"
	"mlauth/pkg/dao"
)

type UserServiceImpl struct {
	UnimplementedUserServiceServer
}

func (s UserServiceImpl) GetUser(_ context.Context, req *GetUserReq) (*GetUserRes, error) {
	u, err := dao.SelectUser(int(req.Uid))
	if err != nil {
		return nil, err
	}

	return &GetUserRes{
		Uid:         int32(u.Uid),
		Email:       u.Email,
		DisplayName: u.DisplayName,
	}, nil
}
