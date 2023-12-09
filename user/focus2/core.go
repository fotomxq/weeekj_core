package UserFocus2

import "errors"

//用户关注模块

// checkMark 检查关注行为特征码
func checkMark(mark string) (err error) {
	switch mark {
	case "focus":
	case "like":
	case "collection":
	default:
		err = errors.New("mark not support")
	}
	return
}

// checkSystem 检查系统来源
func checkSystem(system string) (err error) {
	switch system {
	case "blog_content":
	case "mall_product":
	case "user_core":
	case "info_exchange":
	default:
		err = errors.New("system not support")
	}
	return
}
