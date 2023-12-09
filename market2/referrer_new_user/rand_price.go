package Market2ReferrerNewUser

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"strings"
)

func getMarketReferrerNewUserDepositPriceFix(configVal string) (fixPrice int64) {
	//初始化
	fixPrice = -1
	//检查值是否正确可用
	if configVal == "" {
		return
	}
	//拆分后临时结构
	type tmpStruct struct {
		//最小金额
		Min int64
		//最大金额
		Max int64
		//权重
		Weight int
	}
	var tmpConfig []tmpStruct
	//权重序列
	var weightArr []int
	//分割配置
	configValArr := strings.Split(configVal, ";")
	for _, v := range configValArr {
		vArr := strings.Split(v, ",")
		if len(vArr) != 2 {
			continue
		}
		v2Arr := strings.Split(vArr[0], "-")
		if len(v2Arr) != 2 {
			continue
		}
		vTmp := tmpStruct{
			Min:    CoreFilter.GetInt64ByStringNoErr(v2Arr[0]),
			Max:    CoreFilter.GetInt64ByStringNoErr(v2Arr[1]),
			Weight: int(CoreFilter.GetFloat64ByStringNoErr(vArr[1]) * 100),
		}
		weightArr = append(weightArr, vTmp.Weight)
		tmpConfig = append(tmpConfig, vTmp)
	}
	//判断是否存在配置
	if len(tmpConfig) < 1 {
		return
	}
	//生成随机数，判断权重
	weightKey := CoreFilter.RandomWeightedValue(weightArr)
	for _, v := range tmpConfig {
		if v.Weight != weightArr[weightKey] {
			continue
		}
		fixPrice = int64(CoreFilter.GetRandNumber(int(v.Min), int(v.Max)))
	}
	//没有找到对应权重，直接按照第一个抽取
	fixPrice = int64(CoreFilter.GetRandNumber(int(tmpConfig[0].Min), int(tmpConfig[0].Max)))
	return
}
