package RouterUserCore

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
)

//服务主体通用的用户部分结构

// GetDataFromByUser 获取CoreFieldsFrom的用户来源
func GetDataFromByUser(userData *UserCore.FieldsUserType) CoreSQLFrom.FieldsFrom {
	return CoreSQLFrom.FieldsFrom{
		System: "user",
		ID:     userData.ID,
		Mark:   "",
		Name:   userData.Name,
	}
}
