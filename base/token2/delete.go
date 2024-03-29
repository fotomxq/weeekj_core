package BaseToken2

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// DeleteToken 清理token
func DeleteToken(id int64) {
	//删除会话
	_, _ = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_token2", "id", map[string]interface{}{
		"id": id,
	})
	deleteTokenCache(id)
	//删除短地址
	_ = DeleteTokenSByTokenID(id)
}

// DeleteTokenByLoginFrom 剔除指定登录渠道的所有会话
func DeleteTokenByLoginFrom(loginFrom string) {
	var page int64 = 1
	for {
		dataList, _, _ := GetList(&ArgsGetList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: page,
				Max:  1000,
				Sort: "id",
				Desc: false,
			},
			UserID:    -1,
			OrgID:     -1,
			OrgBindID: -1,
			DeviceID:  -1,
			LoginFrom: loginFrom,
			Search:    "",
		})
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			DeleteToken(v.ID)
		}
		page += 1
	}
}
