package BaseEmail2

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

var (
	// baseEmail2SQL 反馈SQL
	baseEmail2SQL CoreSQL2.Client
)

func Init() {
	baseEmail2SQL.Init(&Router2SystemConfig.MainSQL, "core_email_template")
}
