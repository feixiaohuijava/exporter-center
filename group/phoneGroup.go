package group

import (
	"exporter-center/handlers/auth"
	"exporter-center/handlers/phoneGroup"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PhoneGroup(phoneRouter *gin.RouterGroup, dbConnection *gorm.DB) {
	phoneRouter.POST("/callphone", auth.JWTAuthMiddleware(), func(context *gin.Context) {
		phoneGroup.CallPhone(context, dbConnection)
	})
}
