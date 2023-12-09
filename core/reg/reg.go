package CoreReg

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"github.com/denisbrodbeck/machineid"
	"strings"
)

//本地化注册机服务
// 本服务致力于解决软件或系统授权问题
// 提供标准的验证和算法功能，作为简单的验证工具
// 验证算法的计算方法：本机器特征码 + 当前年和月 + APP特征序列 + 部分特征序列抽取
// 注意，由于当前年和月，可能会产生一系列结构，所以将自动计算出所有可能的序列号，只要用户给予的key满足，即验证通过

var (
	//SNApp 应用特征序列
	// 例如：app1.0
	SNApp string
)

// Init 初始化
func Init(snApp string) {
	SNApp = snApp
}

// Verify 验证序列号是否匹配
func Verify(key string) bool {
	//从key中抽取授权时间
	keys := strings.Split(key, "-")
	if len(keys) != 2 || len(keys[1]) != 12 {
		return false
	}
	//将开始和结束时间抽取出来
	startTime := string(keys[1][0]) + string(keys[1][1]) + string(keys[1][2]) + string(keys[1][3]) + string(keys[1][4]) + string(keys[1][5])
	endTime := string(keys[1][6]) + string(keys[1][7]) + string(keys[1][8]) + string(keys[1][9]) + string(keys[1][10]) + string(keys[1][11])
	//计算真正的key
	code, err := GetCode()
	if err != nil {
		return false
	}
	keyTrue := GetKey(code, startTime, endTime)
	//如果不正确，反馈
	if keyTrue != key {
		return false
	}
	//检查时间是否符合当前时间
	startTimeInt, err := CoreFilter.GetIntByString(startTime)
	if err != nil {
		return false
	}
	endTimeInt, err := CoreFilter.GetIntByString(endTime)
	if err != nil {
		return false
	}
	nowTimeStr := CoreFilter.GetNowTime().Format("200601")
	nowTimeInt, err := CoreFilter.GetIntByString(nowTimeStr)
	if err != nil {
		return false
	}
	if nowTimeInt < startTimeInt {
		return false
	}
	if nowTimeInt > endTimeInt {
		return false
	}
	//全部通过则反馈真
	return true
}

// GetCode 获取本机序列号
func GetCode() (string, error) {
	localCode, err := getLocalCode()
	if err != nil {
		return "", err
	}
	newCode := string(localCode[2]) + string(localCode[6]) + string(localCode[31]) + string(localCode[22]) + string(localCode[17]) + string(localCode[9]) + string(localCode[11]) + string(localCode[13]) + string(localCode[16]) + string(localCode[14]) + string(localCode[20]) + string(localCode[22]) + string(localCode[33]) + string(localCode[24]) + string(localCode[21]) + string(localCode[36]) + string(localCode[14]) + string(localCode[34]) + string(localCode[12]) + string(localCode[1])
	return newCode, nil
}

// GetKey 根据code计算获取key结果
// param code string 机器代码
// param startTime string 开始年月
// param endTime string 结束年月
// return string 计算结果
func GetKey(code string, startTime string, endTime string) string {
	return GetKeyAndApp(SNApp, code, startTime, endTime)
}

// GetKeyAndApp 带有版本序列的生成工具
func GetKeyAndApp(app string, code string, startTime string, endTime string) string {
	newCode := CoreFilter.GetSha1Str(code + app + startTime + endTime)
	key := string(newCode[31]) + string(newCode[1]) + string(newCode[6]) + string(newCode[17]) + string(newCode[7]) + string(newCode[9]) + string(newCode[11]) + string(newCode[13]) + string(newCode[16]) + string(newCode[12]) + string(newCode[20]) + string(newCode[22]) + string(newCode[27]) + string(newCode[25]) + string(newCode[30]) + string(newCode[36]) + string(newCode[37]) + string(newCode[33]) + string(newCode[29]) + string(newCode[14]) + "-" + startTime + endTime
	return key
}

// getLocalCode 获取本机特征码
func getLocalCode() (string, error) {
	id, err := machineid.ProtectedID(SNApp)
	if err != nil {
		return "", err
	}
	return id, nil
}

// getSHA1 获取字符串的sha1
func getSHA1(str string) (string, error) {
	return CoreFilter.GetSha1ByString(str)
}
