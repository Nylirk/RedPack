package utils

import (
	"fmt"
	"testing"
)

func TestSnowflakeGenerate(t *testing.T) {
	workerID := int64(1)
	datacenterID := int64(1)

	sf, err := NewSnowflake(workerID, datacenterID)
	if err != nil {
		t.Fatalf("创建Snowflake实例失败: %v", err)
	}

	ids := make(map[int64]struct{})

	for i := 0; i < 10; i++ {
		id, err := sf.Generate()
		if err != nil {
			t.Errorf("生成ID失败: %v", err)
			continue
		}
		if _, exists := ids[id]; exists {
			t.Errorf("生成的ID重复: %d", id)
		}
		ids[id] = struct{}{}
		fmt.Printf("生成的ID: %d\n", id)
	}
}
