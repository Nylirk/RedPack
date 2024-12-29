package system

import (
	"RedPack/model/common/response"
	"RedPack/model/system"
	"RedPack/model/system/request"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type RedPackApi struct {
}

func (ra *RedPackApi) CreateRedPack(c *gin.Context) {
	var r request.CreateRedPackRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	redPack := &system.RedPack{
		SurplusAmount: r.Amount,
		SurplusTotal:  r.Total,
		TotalAmount:   r.Amount,
		Total:         r.Total,
		UserID:        uuid.Must(uuid.NewV4()),
	}
	id, err := redPackService.CreateRedPackService(redPack)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(id, "创建成功", c)
}

func (ra *RedPackApi) GetRedPack(c *gin.Context) {
	var r request.GetRedPackRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, amount := redPackService.GetRedPackService(r.RedPackID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(amount, "", c)
}

func (ra *RedPackApi) ViewRedPack(c *gin.Context) {
	var r request.GetRedPackRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, IDs := redPackService.ViewRedPackService(r.RedPackID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(IDs, "", c)
}
