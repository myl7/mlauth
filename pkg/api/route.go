package api

import "github.com/gin-gonic/gin"

func Run() error {
	r := gin.Default()
	api := r.Group("/api")

	api.POST("/users/login", userLogin)
	api.POST("/users/renew", userRenew)
	api.POST("/users", userRegister)

	api.POST("/emails/active", emailActive)
	api.POST("/emails/active/retry", userAuthExist, emailActiveRetry)

	authed := api.Group("/")
	authed.Use(userAuth)
	authed.GET("/users/me", userGet)
	authed.PUT("/users/me", userEdit)

	err := r.Run()
	if err != nil {
		return err
	}

	return nil
}
