package system_test

import (
	"RedPack/global"
	"RedPack/model/system"
	service "RedPack/service/system"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"testing"
)

const (
	testRedPackUID = "test_red_pack_uid"
	numGoroutines  = 100
)

// 初始化测试数据库
func TestMain(m *testing.M) {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/redpack?charset=utf8mb4&parseTime=True&loc=Local"
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移模式
	global.DB.AutoMigrate(&system.RedPack{}, &system.RedPackRecord{})

	// 初始化一些测试数据
	testRedPack := &system.RedPack{
		RedPackUID:    testRedPackUID,
		SurplusAmount: 100,
		SurplusTotal:  10,
	}
	if err := global.DB.Create(testRedPack).Error; err != nil {
		panic("failed to create test red pack")
	}

	m.Run()
}

// 并发测试 GetRedPackService
func TestGetRedPackServiceConcurrent(t *testing.T) {
	rs := service.RedPackService{}
	redPackUID := testRedPackUID

	// 初始化红包
	testRedPack := &system.RedPack{
		RedPackUID:    redPackUID,
		SurplusAmount: 100,
		SurplusTotal:  10,
	}
	if err := global.DB.Create(testRedPack).Error; err != nil {
		t.Fatalf("failed to create test red pack: %v", err)
	}

	var wg sync.WaitGroup
	results := make(chan float64, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err, amount := rs.GetRedPackService(redPackUID)
			if err != nil {
				errors <- err
			} else {
				results <- amount
			}
		}()
	}

	wg.Wait()
	close(results)
	close(errors)

	// 检查错误
	for err := range errors {
		t.Errorf("concurrent request failed: %v", err)
	}

	// 检查结果
	totalAmount := 0.0
	for amount := range results {
		totalAmount += amount
	}

	// 验证总金额是否正确
	expectedTotalAmount := 100.0
	if totalAmount != expectedTotalAmount {
		t.Errorf("expected total amount %v, but got %v", expectedTotalAmount, totalAmount)
	}
}
