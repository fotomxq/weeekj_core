package OrgSign

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

//组织成员签名
/**
用途：
1. 用于存储和传递组织成员的签名
2. 也可以用于用户的签名
*/

var (
	//签名存储
	signDB BaseSQLTools.Quick
)

func Init() (err error) {
	//初始化指标定义
	if err = signDB.Init("org_sign", &FieldsSign{}); err != nil {
		return
	}
	//反馈
	return
}
