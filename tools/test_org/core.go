package TestOrg

import (
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	OrgTime "github.com/fotomxq/weeekj_core/v5/org/time"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

var (
	UserInfo     UserCore.FieldsUserType
	OrgData      OrgCore.FieldsOrg
	WorkTimeData OrgTime.FieldsWorkTime
	GroupData    OrgCore.FieldsGroup
	BindData     OrgCore.FieldsBind
	BindList     []OrgCore.FieldsBind
)
