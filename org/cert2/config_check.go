package OrgCert2

import "errors"

// 检查审核类型
func checkAuditType(t string) (err error) {
	switch t {
	case "none":
	case "wait":
	case "auto":
	default:
		err = errors.New("unknown audit type")
		return
	}
	return
}

// 检查提醒类型
func checkTipType(t string) (err error) {
	switch t {
	case "none":
	case "audit":
	case "expire":
	case "all":
	default:
		err = errors.New("unknown audit type")
		return
	}
	return
}
