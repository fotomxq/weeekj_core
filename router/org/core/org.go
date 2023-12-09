package RouterOrgCore

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
)

//服务主体通用的组织结构部分

// 获取CoreFieldsFrom的组织来源
// 用于创建和更新过程
func GetDataFromByOrg(orgData *OrgCore.FieldsOrg) CoreSQLFrom.FieldsFrom {
	return CoreSQLFrom.FieldsFrom{
		System: "org",
		ID:     orgData.ID,
		Mark:   "",
		Name:   orgData.Name,
	}
}

// 无名称指定
// 用于常用搜索事项
func GetDataFromByOrgNoName(orgData *OrgCore.FieldsOrg) CoreSQLFrom.FieldsFrom {
	return CoreSQLFrom.FieldsFrom{
		System: "org",
		ID:     orgData.ID,
		Mark:   "",
		Name:   "",
	}
}
