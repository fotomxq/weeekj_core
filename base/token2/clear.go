package BaseToken2

import CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"

// 启动服务后，自动清理一次过期的会话
func clearExpireToken() {
	//遍历
	var page int64 = 1
	for {
		//获取过期的会话
		var dataList []FieldsToken
		_ = baseToken.Select().SetPages(CoreSQL2.ArgsPages{
			Page: page,
			Max:  1000,
			Sort: "id",
			Desc: false,
		}).SetFieldsList([]string{"id"}).SelectList("expire_at < NOW()").Result(&dataList)
		//检查内容
		if len(dataList) < 1 {
			break
		}
		//遍历数据
		for _, v := range dataList {
			DeleteToken(v.ID)
		}
		//下一页
		page += 1
	}
}
