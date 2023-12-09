package CoreFilter

// CheckDes 检查描述信息
func CheckDes(str string, min int, max int) bool {
	strLen := len([]rune(str))
	if strLen < min || strLen > max {
		return false
	}
	return true
}

// CheckContent 检查消息内容
func CheckContent(str string, min int, max int) bool {
	strLen := len([]rune(str))
	if strLen < min || strLen > max {
		return false
	}
	return true
}

// CheckEmail 验证邮箱
// param str string 邮箱地址
// return bool 是否正确
func CheckEmail(str string) bool {
	if str == "" {
		return false
	}
	return MatchStr(`^([\w\.\_]{1,225})@(\w{1,}).(([a-z]{1,225})|([a-z]{1,225}).([a-z]{1,225}))$`, str)
}

// CheckExpireTime 检查过期时间
func CheckExpireTime(expireTime string) bool {
	return MatchStr(`^[0-9a-z-]{1,10}$`, expireTime)
}

// CheckFilterStr 过滤非法字符后判断其长度是否符合标准
// param str string 要过滤的字符串
// param min int 最短，包括该长度
// param max int 最长，包括该长度
// return string 过滤后的字符串，失败返回空字符串
func CheckFilterStr(str string, min int, max int) string {
	newStr := FilterStrForce(str)
	if newStr == "" {
		return ""
	}
	strLen := len(newStr)
	if strLen >= min && strLen <= max {
		return newStr
	}
	return ""
}

// CheckFromNameAndEmpty 检查可以为空的字段名称信息
func CheckFromNameAndEmpty(str string) bool {
	if str == "" {
		return true
	}
	return MatchStr(`^[a-zA-Z0-9_-]{1,20}$`, str)
}

// CheckHexSha1 验证是否为SHA1
// param str string 字符串
// return bool 是否正确
func CheckHexSha1(str string) bool {
	return MatchStr(`^[a-zA-Z0-9]{40}$`, str)
}

// CheckHexSha256 验证是否为SHA256
// param str string 字符串
// return bool 是否正确
func CheckHexSha256(str string) bool {
	return MatchStr(`^[a-zA-Z0-9]{64}$`, str)
}

// CheckID 检查ID
// param str string ID序列
// return bool 是否正确
func CheckID(str string) bool {
	id, err := GetInt64ByString(str)
	if err != nil {
		return false
	}
	if id < 1 {
		return false
	}
	return true
}

// CheckSN 检查SN
func CheckSN(sn int64) bool {
	if sn > 0 {
		return true
	}
	return false
}

// CheckIDCard 验证身份证
// 因为复杂性，仅考虑验证身份证位数有效性
// 未来可根据实际需求加入外部API对身份证进行二次验证
// param str string 身份证号码
// return bool 是否正确
func CheckIDCard(str string) bool {
	if str == "" {
		return false
	}
	if len(str) > 10 && len(str) < 20 {
		return true
	}
	return false
}

// CheckIP 验证是否为IP地址
// param str string IP地址
// return bool 是否正确
func CheckIP(str string) bool {
	if str == "[::1]" || str == "::1" || str == "localhost" {
		return true
	}
	if MatchStr(`((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`, str) {
		return true
	}
	if MatchStr(`^$(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`, str) {
		return true
	}
	return false
}

// CheckHost 验证host
func CheckHost(str string) bool {
	if CheckIP(str) {
		return true
	}
	if MatchStr(`(http://|https://)?([^/]*)`, str) {
		return true
	}
	return false
}

// CheckPort 验证port
func CheckPort(str string) bool {
	return MatchStr(`^[0-9]{1,10}`, str)
}

// CheckLimit 验证传统限制
func CheckLimit(limit int, min int, max int) bool {
	if limit >= min && limit <= max {
		return true
	}
	return false
}

// CheckNiceName 检查昵称
// param str string 昵称
// return bool 是否正确
func CheckNiceName(str string) bool {
	if str == "" {
		return false
	}
	if len(str) < 0 || len(str) > 50 {
		return false
	}
	return true
}

// CheckCityCode 检查地区编码
func CheckCityCode(str string) bool {
	return MatchStr(`^[a-zA-Z0-9_-]{2,30}$`, str)
}

// CheckAddress 检查详细收货地址
func CheckAddress(str string) bool {
	if len(str) < 5 || len(str) > 300 {
		return false
	}
	return true
}

// CheckNationCode 检查手机号国家代码
// 86
func CheckNationCode(str string) bool {
	return MatchStr(`^[0-9]{1,6}$`, str)
}

// CheckMark 检查mark
// param str string
// return bool 是否正确
func CheckMark(str string) bool {
	if str == "" {
		return false
	}
	return MatchStr(`^[a-zA-Z0-9_-]{1,100}$`, str)
}

func CheckMarkPage(str string) bool {
	if str == "" {
		return false
	}
	return MatchStr(`\/([^?#]{1,50})`, str)
}

// CheckTimeType 检查时间类型
func CheckTimeType(str string) bool {
	switch str {
	case "year":
	case "month":
	case "day":
	case "hour":
	case "minute":
	case "second":
	default:
		return false
	}
	return true
}

// CheckPassword 验证密码
// param str string 密码
// return bool 是否正确
func CheckPassword(str string) bool {
	if str == "" {
		return false
	}
	return MatchStr(`^[a-zA-Z0-9_()\[\]-]{6,30}$`, str)
}

// CheckPhone 验证电话号码
// 必须是手机电话号码或带区号的固定电话号码
// eg 03513168322
// eg 13066889999
func CheckPhone(str string) bool {
	if str == "" {
		return false
	}
	return MatchStr(`^[0-9]{11}$`, str)
}

// CheckPage 检查页数
func CheckPage(page int64) bool {
	return page > 0
}

// CheckMax 检查页长
func CheckMax(max int64) bool {
	return max < 1000
}

// CheckSearch 验证搜索类型的字符串
// param str string 字符串
// return bool 是否正确
func CheckSearch(str string) bool {
	if str == "" {
		return false
	}
	return MatchStr(`[\p{Han}|a-zA-Z0-9\\\^\.\,\*\+\?\{\}\(\)\[\]\|]{1,300}`, str)
}

// CheckSort 检查字段名称
func CheckSort(sort string) bool {
	return MatchStr(`^[a-zA-Z0-9_-]{1,20}$`, sort)
}

// CheckUsername 检查用户名
// param str string 用户名
// return bool 是否正确
func CheckUsername(str string) bool {
	if str == "" {
		return false
	}
	return MatchStr(`^[a-zA-Z0-9_-]{4,20}$`, str)
}

// CheckVcode 检查验证码
func CheckVcode(str string) bool {
	if str == "" {
		return false
	}
	return MatchStr(`^[a-zA-Z0-9_-]{3,30}$`, str)
}

// CheckFileName 检查文件名称
func CheckFileName(str string) bool {
	if str == "" {
		return false
	}
	//return MatchStr(`^[\u4e00-\u9fa5_a-zA-Z0-9.]+$`, str)
	return MatchStr(`^[\p{Han}|a-zA-Z0-9\.]{3,100}`, str)
}

// CheckMapType 检查GPS标准
func CheckMapType(str int) bool {
	switch str {
	case 0:
	//case "WGS-84":
	case 1:
	//case "GCJ-02":
	case 2:
	//case "BD-09":
	default:
		return false
	}
	return true
}

// CheckGPS 检查GPS坐标
func CheckGPS(str float64) bool {
	if str > -255 || str < 255 {
		return true
	}
	return false
}
