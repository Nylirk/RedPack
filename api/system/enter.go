package system

import "RedPack/service"

type ApiGroup struct {
	RedPackApi
}

var (
	redPackService = service.ServiceGroupApp.SystemServiceGroup.RedPackService
)
