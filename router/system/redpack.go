package system

import (
	"RedPack/api"
	"github.com/gin-gonic/gin"
)

type RedPackRouter struct {
}

func (r *RedPackRouter) InitRedPackRouter(Router *gin.RouterGroup) {
	router := Router.Group("rp")
	redPackApi := api.ApiGroupApp.SystemApiGroup.RedPackApi
	{
		router.POST("create", redPackApi.CreateRedPack)
		router.POST("get", redPackApi.GetRedPack)
		router.POST("view", redPackApi.ViewRedPack)
	}
}
