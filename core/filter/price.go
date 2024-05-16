package CoreFilter

import "fmt"

// GetPriceByUint8 转化uint8为int64结构
func GetPriceByUint8(data []uint8) int64 {
	newDataF, _ := GetFloat64ByUint8(data)
	newData := GetInt64ByFloat64(newDataF * 100)
	return newData
}

// GetPriceToShowPrice 将金额转为float64并保留2位，显示使用的金额
func GetPriceToShowPrice(data int64) float64 {
	return float64(data) / 100
}

// GetPriceToShowPriceStr 将金额转为float64并保留2位，显示使用的金额
func GetPriceToShowPriceStr(data float64) string {
	return fmt.Sprintf("%.2f", data)
}
