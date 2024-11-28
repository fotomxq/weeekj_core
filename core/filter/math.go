package CoreFilter

import (
	"math"
)

//数学扩展工具

// GetRound 四舍五入取整数
func GetRound(data float64) float64 {
	if math.IsNaN(data) || math.IsInf(data, 0) {
		return 0
	}
	return math.Floor(data) + 0.5
}

// RoundToTwoDecimalPlaces 四舍五入保留2位小数点
func RoundToTwoDecimalPlaces(num float64) float64 {
	if math.IsNaN(num) || math.IsInf(num, 0) {
		return 0
	}
	return math.Round(num*100) / 100 // 四舍五入到2位整数并返回float64类型结果
}

// RoundTo4DecimalPlaces 四舍五入保留4位小数点
func RoundTo4DecimalPlaces(num float64) float64 {
	if math.IsNaN(num) || math.IsInf(num, 0) {
		return 0
	}
	return math.Round(num*10000) / 10000 // 四舍五入到2位整数并返回float64类型结果
}

func GetRoundToInt(data float64) int {
	return int(math.Floor(data) + 0.5)
}

func GetRoundToInt64(data float64) int64 {
	return int64(math.Floor(data) + 0.5)
}

// MathLastProportion 计算提升后占比
func MathLastProportion(prev int64, last int64) (addCount int64, p float64) {
	addCount = last - prev
	if prev == 0 && last > 0 {
		p = 1
		return
	}
	if addCount == 0 {
		p = 0
		return
	}
	if prev >= 0 && last == 0 {
		p = -1
		return
	}
	p = (float64(last) / float64(prev)) - 1
	return
}

func MathLastProportionToInt64(prev int64, last int64) (addCount int64, p int64) {
	var p2 float64
	addCount, p2 = MathLastProportion(prev, last)
	p = int64(p2 * 10000)
	return
}
