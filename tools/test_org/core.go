package TestOrg

import (
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	OrgTime "gitee.com/weeekj/weeekj_core/v5/org/time"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
)

var (
	UserInfo     UserCore.FieldsUserType
	OrgData      OrgCore.FieldsOrg
	WorkTimeData OrgTime.FieldsWorkTime
	GroupData    OrgCore.FieldsGroup
	BindData     OrgCore.FieldsBind
	BindList     []OrgCore.FieldsBind
)
