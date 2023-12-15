package AnalysisAny2

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 删除指定数据
func deleteAnyByID(id int64) {
	_, err := CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "analysis_any2", "id", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return
	}
}
