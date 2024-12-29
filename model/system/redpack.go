package system

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type RedPack struct {
	ID            int64 `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time
	UserID        uuid.UUID `json:"userId" gorm:"index;comment:用户UUID"`
	TotalAmount   float64   `json:"totalAmount" gorm:"comment:总金额"`
	SurplusAmount float64   `json:"surplusAmount" gorm:"comment:剩余金额"`
	Total         int       `json:"total" gorm:"comment:红包总数"`
	SurplusTotal  int       `json:"surplusTotal" gorm:"comment:剩余红包总数"`
}

// TableName 返回表名
func (RedPack) TableName() string {
	return "red_pack"
}
