package initialize

import (
	"RedPack/middleware"
	"RedPack/router"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	Router := gin.New()
	Router.Use(middleware.Cors())
	systemRouter := router.RouterGroupApp.System
	systemGroup := Router.Group("")
	systemRouter.InitRedPackRouter(systemGroup)
	return Router
}
