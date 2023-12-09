package BlogStuRead

import "errors"

// 检查文章类型
func checkContentType(s int) (err error) {
	switch s {
	case 0:
	case 1:
	case 2:
	default:
		err = errors.New("content type not support")
	}
	return
}
