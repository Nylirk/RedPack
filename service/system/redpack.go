package system

import (
	"RedPack/global"
	"RedPack/model/system"
	"RedPack/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"math/rand"
	"strconv"
	"time"
)

type RedPackService struct {
}

func (rs *RedPackService) CreateRedPackService(r *system.RedPack) (err error) {
	var rp system.RedPack
	// 创建分布式锁实例
	lock := utils.NewRedisDistributedLock(global.REDIS, "lock_key", "client_id", 5*time.Second)
	// 尝试获取锁
	if lock.TryLock() {
		defer func() error {
			if err = lock.Unlock(); err != nil {
				return errors.New("解锁失败")
			}
			return nil
		}()
		workerID := int64(1)
		datacenterID := int64(1)
		sf, err := utils.NewSnowflake(workerID, datacenterID)
		if err != nil {
			return errors.New("创建Snowflake实例失败")
		}
		id, err := sf.Generate()
		if err != nil {
			return errors.New("生成ID失败")
		}
		r.ID = id
	}

	tx := global.DB.Clauses(dbresolver.Write).Begin()
	defer func() {
		if re := recover(); re != nil {
			tx.Rollback()
		}
	}()
	if !errors.Is(tx.Where("id= ?", r.ID).First(&rp).Error, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return errors.New("红包已创建")
	}
	if err = tx.Create(&r).Error; err != nil {
		tx.Rollback()
		return errors.New("创建红包失败")
	}
	err = tx.Commit().Error
	return err
}

func (rs *RedPackService) GetRedPackService(redPackID int64) (error, float64) {
	var rp system.RedPack
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Where("id = ?", redPackID).First(&rp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("红包不存在"), 0
		}
		tx.Rollback()
		return errors.New("查询红包失败"), 0
	}
	if rp.SurplusAmount == 0 || rp.SurplusTotal == 0 {
		return errors.New("红包已被抢完"), 0
	}
	amount, surplusTotal, surplusAmount := utils.GenerateRedPack(rp.SurplusAmount, rp.SurplusTotal)
	newRp := make(map[string]interface{})
	newRp["surplus_amount"] = surplusAmount
	newRp["surplus_total"] = surplusTotal
	if err := tx.Model(&system.RedPack{}).Where("id = ?", redPackID).
		Updates(map[string]interface{}{
			"surplus_amount": surplusAmount,
			"surplus_total":  surplusTotal,
		}).Error; err != nil {
		tx.Rollback()
		return errors.New("更新失败"), 0
	}
	rr := &system.RedPackRecord{
		Amount:    amount,
		RedPack:   rp,
		RedPackId: rp.ID,
		UserID:    uuid.Must(uuid.NewV4()),
	}
	if err := tx.Create(&rr).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建抢红包记录失败: %v", err), 0
	}
	err := tx.Commit().Error
	return err, amount
}

func (rs *RedPackService) ViewRedPackService(redPackID int64) (err error, usersIDs interface{}) {
	var rp system.RedPack
	type RedPackInfo struct {
		UserID uuid.UUID `gorm:"column:user_id"`
		Amount float64   `gorm:"column:amount"`
		Time   string    `gorm:"column:created_at"`
	}
	var redPackInfo []RedPackInfo
	val, err := global.REDIS.Get(context.Background(), strconv.FormatInt(redPackID, 10)).Result()
	if errors.Is(err, redis.Nil) || val == "" {
		tx := global.DB.Begin()
		defer func() {
			if re := recover(); re != nil {
				tx.Rollback()
			}
		}()
		if errors.Is(tx.Where("id = ?", redPackID).First(&rp).Error, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("红包不存在"), nil
		}
		if err = tx.Model(&system.RedPackRecord{}).Where("red_pack_id = ?", rp.ID).Select("user_id, amount, created_at").Find(&redPackInfo).Error; err != nil {
			tx.Rollback()
			return errors.New("未查询到记录"), nil
		}
		jsonByte, err := json.Marshal(redPackInfo)
		if err != nil {
			return errors.New("json转换失败"), nil
		}
		// 设置随机过期时间防止缓存雪崩
		expiration := time.Minute * time.Duration(rand.Intn(4))
		err = global.REDIS.Set(context.Background(), strconv.FormatInt(redPackID, 10), string(jsonByte), expiration).Err()
		if err != nil {
			return errors.New("redis写入失败"), nil
		}
		return tx.Commit().Error, redPackInfo
	}
	err = json.Unmarshal([]byte(val), &usersIDs)
	return err, usersIDs
}
