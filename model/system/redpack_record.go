package system

import (
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"time"
)

type RedPackRecord struct {
	ID        int64 `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	RedPackId int64          `json:"redPackId" gorm:"index;not null;foreignKey:ID;references:ID"`
	UserID    uuid.UUID      `json:"userId" gorm:"index;not null"`
	Amount    float64        `json:"amount"`
	RedPack   RedPack        `gorm:"foreignKey:RedPackId;references:ID"`
}

// TableName 返回表名
func (RedPackRecord) TableName() string {
	return "red_pack_record"
}
