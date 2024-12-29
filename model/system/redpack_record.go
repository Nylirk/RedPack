package system

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type RedPackRecord struct {
	ID        int64 `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	RedPackId int64     `json:"redPackId" gorm:"index;not null;foreignKey:ID;references:ID"`
	UserID    uuid.UUID `json:"userId" gorm:"index;not null"`
	Amount    float64   `json:"amount"`
}

// TableName 返回表名
func (RedPackRecord) TableName() string {
	return "red_pack_record"
}
