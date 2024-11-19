package AnalysisIndexRawFilter

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"regexp"
	"strconv"
	"strings"
)

// FilterPriceAuto 全自动处理模式
// 仅用于字符串类型的数据，如果是其他类型数据，需转化后提供
// 1. 剔除无效的数据
// 2. 如果数据中发现“万”，则乘以10000
// 3. 四舍五入保留两位小数
func FilterPriceAuto(price string) (result float64) {
	//转换
	result = FilterPrice(price)
	if strings.Contains(price, "万") {
		result = result * 10000
	}
	result = CoreFilter.RoundToTwoDecimalPlaces(result)
	//返回
	return
}

// FilterPrice 处理金额、数值类数据
func FilterPrice(price any) (result float64) {
	//转换
	switch price.(type) {
	case float64:
		result = price.(float64)
	case int:
		result = float64(price.(int))
	case int64:
		result = float64(price.(int64))
	case string:
		//尝试转换
		// 通过正则表达式regexp，提取开头为数字的内容
		noComma := regexp.MustCompile(`,`).ReplaceAllString(price.(string), ``)
		re := regexp.MustCompile(`^\d+(\.\d+)?`)
		match := re.FindString(noComma)
		result, _ = strconv.ParseFloat(match, 64)
	}
	//返回
	return
}
