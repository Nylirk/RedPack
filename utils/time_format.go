package utils

import (
	"fmt"
	"time"
)

func TimeFormat(timeStr string) string {
	// 定义原始时间和目标时间的布局
	layoutISO := "2006-01-02T15:04:05.99+07:00"
	layoutTarget := "2006年 01月02日 15:04"

	t, err := time.Parse(layoutISO, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}
	t = t.In(time.FixedZone("CST", 8*60*60))
	formattedTime := t.Format(layoutTarget)

	return formattedTime
}
