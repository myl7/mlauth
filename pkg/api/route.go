package api

import "github.com/gin-gonic/gin"

func Run() error {
	r := gin.Default()
	api := r.Group("/api")

	api.POST("/users/login", userLogin)
	api.POST("/users/renew", userRenew)

	authed := api.Group("/")
	authed.Use(userAuth)

	err := r.Run()
	if err != nil {
		return err
	}

	return nil
}
