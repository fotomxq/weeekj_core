package ServiceAD

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 清理过期投放绑定
func runEnd() {
	//删除过期数据
	_, err := CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_ad_bind", "delete_at < to_timestamp(1000000) AND end_at < NOW()", nil)
	if err != nil {
		//不记录数据
	}
}
