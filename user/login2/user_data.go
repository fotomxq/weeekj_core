package UserLogin2

import (
	BaseFileSys2 "github.com/fotomxq/weeekj_core/v5/base/filesys2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

// DataGetUserData 登录时的数据包汇总
type DataGetUserData struct {
	//用户ID
	ID int64 `json:"id"`
	//用户昵称
	Name string `json:"name"`
	//用户头像
	Avatar string `json:"avatar"`
	//用户权限列
	Permissions []string `json:"permissions"`
}

func GetUserData(userID int64) (data DataGetUserData) {
	//获取用户数据
	userInfo, _ := UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    userID,
		OrgID: -1,
	})
	if userInfo.ID < 1 {
		return
	}
	//组装数据
	data = DataGetUserData{
		ID:          userInfo.ID,
		Name:        userInfo.Name,
		Avatar:      BaseFileSys2.GetPublicURLByClaimID(userInfo.Avatar),
		Permissions: GetUserPermissionsByUserInfo(&userInfo),
	}
	//反馈
	return
}

// GetUserPermissions 获取用户登录常用的权限组
func GetUserPermissions(userID int64) []string {
	//获取用户数据
	userInfo, _ := UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    userID,
		OrgID: -1,
	})
	if userInfo.ID < 1 {
		return []string{}
	}
	//反馈权限
	return GetUserPermissionsByUserInfo(&userInfo)
}

func GetUserPermissionsByUserInfo(userInfo *UserCore.FieldsUserType) []string {
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
