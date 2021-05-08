package group

import (
	"exporter-center/handlers/auth"
	"github.com/gin-gonic/gin"
)

func AuthGroup(phoneRouter *gin.RouterGroup) {
	phoneRouter.POST("/login", func(context *gin.Context) {
		auth.AuthHandler(context)
	})
	phoneRouter.GET("/home", auth.JWTAuthMiddleware(), func(context *gin.Context) {
		auth.HomeHandler(context)
	})
}
