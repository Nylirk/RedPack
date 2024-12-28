package utils

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

type RedPack struct {
	SurplusAmount float64 // 剩余金额
	SurplusTotal  uint    // 红包剩余个数
}

// 取两位小数
func remainTwoDecimal(num float64) float64 {
	numStr := strconv.FormatFloat(num, 'f', 2, 64)
	num, _ = strconv.ParseFloat(numStr, 64)
	return num
}

func GenerateRedPack(surplusAmount float64, surplusTotal int) (float64, int, float64) {
	surplusAmount -= 0.01 * float64(surplusTotal)
	if surplusTotal <= 0 {
		return 0, 0, 0
	}

	if surplusTotal == 1 {
		return remainTwoDecimal(surplusAmount + 0.01), 0, 0
	}

	avgAmount := math.Floor(100*(surplusAmount/float64(surplusTotal))) / float64(100)
	avgAmount = remainTwoDecimal(avgAmount)

	rand.NewSource(time.Now().UnixNano())

	var maxNum float64
	if avgAmount > 0 {
		maxNum = 2*avgAmount - 0.01
	} else {
		maxNum = 0
	}
	money := remainTwoDecimal(rand.Float64()*(maxNum) + 0.01)

	surplusTotal -= 1
	surplusAmount = remainTwoDecimal(surplusAmount + 0.01 - money)

	return money, surplusTotal, surplusAmount
}
