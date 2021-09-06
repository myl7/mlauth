package rpc

import (
	"google.golang.org/grpc"
	user "mlauth/pkg/rpc/user"
)

func Register() *grpc.Server {
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &user.UserServiceImpl{})
	return s
}
