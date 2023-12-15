package UserLogin2URL

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

func getUserPermissionsByUserInfo(userInfo *UserCore.FieldsUserType) []string {
	var permissions []string
	var groupIDs []int64
	for _, v := range userInfo.Groups {
		if v.ExpireAt.Unix() < 1000000 || v.ExpireAt.Unix() >= CoreFilter.GetNowTime().Unix() {
			groupIDs = append(groupIDs, v.GroupID)
		}
	}
	if len(groupIDs) > 0 {
		permissions = UserCore.GetGroupPermissionList(groupIDs)
	}
	return permissions
}
