package CoreFilter

// CheckCountry 检查国家英文标示
func CheckCountry(str int) bool {
	switch str {
	case 86:
	default:
		return false
	}
	return true
}

// CheckProvince 检查省份
func CheckProvince(str int) bool {
	if str > 0 && str < 9999999 {
		return true
	}
	return false
}
