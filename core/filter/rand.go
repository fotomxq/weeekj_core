package CoreFilter

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

// GetRandStr 获取随机字符串
// param n int 随机码
// return string 新随机字符串
func GetRandStr(n int) string {
	r := rand.New(rand.NewSource(GetNowTime().UnixNano()))
	re := r.Intn(n)
	return strconv.Itoa(re)
}

// GetRandStr2 获取随机字符串
// 第二代，自动生成混淆随机数，与时间混淆后计算得出
func GetRandStr2() string {
	rn := rand.Intn(9999999999)
	r := rand.New(rand.NewSource(GetNowTime().UnixNano()))
	re := r.Intn(rn)
	return strconv.Itoa(re)
}

// GetRandStr3 获取随机数
// 第三代，支持长度限制
func GetRandStr3(limit int) (string, error) {
	var str string
	var err error
	if limit < 1 {
		// < 1 抛出空字符串
		return "", errors.New("get rand str, but limit is 0.")
	} else if limit < 40 && limit > 0 {
		// <= 40 && > 0 则采用sha1截取
		randStr := GetRandStr2()
		str = GetSha1Str(randStr)
		if str == "" {
			return "", errors.New("get rand str, but sha1 is un know error.")
		}
		str = string([]byte(str)[:limit])
	} else if limit == 40 {
		//  < 40 则采用sha1截取
		randStr := GetRandStr2()
		str = GetSha1Str(randStr)
		if str == "" {
			return "", errors.New("get rand str, but sha1 is un know error.")
		}
	} else if limit > 40 && limit < 56 {
		// > 40 && < 56 则采用sha256截取
		randStr := GetRandStr2()
		str, err = GetSha256Str(randStr)
		if err != nil {
			return "", err
		}
		str = string([]byte(str)[:limit])
	} else if limit == 56 {
		// == 56 则采用sha256随机
		randStr := GetRandStr2()
		str, err = GetSha256Str(randStr)
		if err != nil {
			return "", err
		}
	} else {
		// > 56 则采用sha256 + 随机数补充
		// 默认将按照sha256 + sha256叠加处理，如果其他特殊需求请自行配置
		randStr1 := GetRandStr2()
		str1, err := GetSha256Str(randStr1)
		if err != nil {
			return "", err
		}
		randStr2 := GetRandStr2()
		str2, err := GetSha256Str(randStr2)
		if err != nil {
			return "", err
		}
		str = str1 + str2
	}
	return str, nil
}

// GetRandStr4 获取随机数
// 隐藏错误信息
func GetRandStr4(limit int) string {
	str, err := GetRandStr3(limit)
	if err != nil {
		str = GetRandStr(limit)
	}
	return str
}

// GetRandNumber 生成指定范围的随机数字
func GetRandNumber(min int, max int) int {
	rand.Seed(GetNowTime().UnixNano())
	res := rand.Intn(max - min)
	res = res + min
	return res
}

// RandomWeightedValue 通过一组int数组作为权重，随机并生成符合条件的权重值key
func RandomWeightedValue(weights []int) (resultKey int) {
	// 确保随机数的种子每次都不同
	rand.Seed(time.Now().UnixNano())
	// 计算权重总和
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}
	// 生成随机数
	randomNum := rand.Intn(totalWeight)
	// 找到随机数所在的权重区间
	weightSum := 0
	for i, weight := range weights {
		weightSum += weight
		if randomNum < weightSum {
			return i
		}
	}
	return -1 // 如果出现错误，返回-1
}
