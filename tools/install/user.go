package ToolsInstall

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"time"
)

func InstallUser() error {
	//配置文件名称
	configFileName := "user_groups.json"
	//检查配置文件是否存在
	if !checkConfigFile(configFileName) {
		return nil
	}
	//是否为追加模式
	installAppend, err := Router2SystemConfig.Cfg.Section("core").Key("install_append").Bool()
	if err != nil {
		installAppend = false
	}
	//加载和安装权限数据，自动重建
	type dataUserPermissionDataType struct {
		Mark string
		Name string
		Des  string
	}
	type dataUserPermissionType struct {
		Permissions []dataUserPermissionDataType
	}
	if !installAppend {
		//不需要删除权限，而是不断追加即可
		//if err := UserCore.DeleteAllPermission(); err != nil {
		//	CoreLog.Info(errors.New("remove all permission, " + err.Error()))
		//	err = nil
		//}
	}
	fileList, err := CoreFile.GetFileList(fmt.Sprint(configDir, "user_permissions"), []string{"json"}, true)
	if err != nil {
		return nil
	}
	for _, vSrc := range fileList {
		var permissionData dataUserPermissionType
		fileByte, err := CoreFile.LoadFile(vSrc)
		if err != nil {
			return errors.New("load file, " + err.Error())
		}
		err = json.Unmarshal(fileByte, &permissionData)
		if err != nil {
			return errors.New("get byte to json, " + err.Error())
		}
		for _, v := range permissionData.Permissions {
			data, err := UserCore.GetPermissionByMark(&UserCore.ArgsGetPermissionByMark{
				Mark: v.Mark,
			})
			if err == nil {
				err = UserCore.DeletePermission(&UserCore.ArgsDeletePermission{
					Mark:           data.Mark,
					SkipCheckGroup: true,
				})
				if err != nil {
					return err
				}
			}
			if err := UserCore.CreatePermission(&UserCore.ArgsCreatePermission{
				Mark:     v.Mark,
				Name:     v.Name,
				Des:      v.Des,
				AllowOrg: false,
			}); err != nil {
				return errors.New("create permission, " + err.Error())
			}
		}
	}
	//加载全部权限
	allPermissions1, err := UserCore.GetAllPermission(&UserCore.ArgsGetAllPermission{
		AllowOrg: false,
	})
	if err != nil {
		return err
	}
	allPermissions2, err := UserCore.GetAllPermission(&UserCore.ArgsGetAllPermission{
		AllowOrg: true,
	})
	if err != nil {
		return err
	}
	var allPermissions []UserCore.FieldsPermissionType
	for _, v := range allPermissions1 {
		allPermissions = append(allPermissions, v)
	}
	for _, v := range allPermissions2 {
		allPermissions = append(allPermissions, v)
	}
	var allPermissionsList []string
	for _, v2 := range allPermissions {
		allPermissionsList = append(allPermissionsList, v2.Mark)
	}
	//tip
	//CoreLog.Info("all permission list: ", allPermissions)
	//尝试获取所有底层用户组
	allGroups, err := UserCore.GetAllGroup(&UserCore.ArgsGetAllGroup{
		OrgID: 0,
	})
	if err != nil {
		err = nil
	}
	//安装用户组数据，仅安装内置的用户组，如果存在则跳过
	type dataUserGroupsDataType struct {
		Name        string
		Des         string
		Permissions []string `json:"Permissions"`
	}
	type dataUserGroupsType struct {
		Groups []dataUserGroupsDataType
	}
	var groupData dataUserGroupsType
	if err := loadConfigFile(configFileName, &groupData); err != nil {
		return nil
	}
	for _, v := range groupData.Groups {
		vPermissions := v.Permissions
		if v.Name == "管理员" {
			vPermissions = allPermissionsList
		}
		//找到名称一致的用户组，将更新，如果没找到则创建新的
		var findData UserCore.FieldsGroupType
		isFindGroup := false
		for _, vGroup := range allGroups {
			if vGroup.Name == v.Name {
				isFindGroup = true
				findData = vGroup
				break
			}
		}
		if isFindGroup {
			//叠加新的权限设计
			for _, v2 := range v.Permissions {
				isFind := false
				for _, v3 := range vPermissions {
					if v2 == v3 {
						isFind = true
						break
					}
				}
				if !isFind {
					vPermissions = append(vPermissions, v2)
				}
			}
			//如果存在，则修正信息
			if err := UserCore.UpdateGroup(&UserCore.ArgsUpdateGroup{
				ID:          findData.ID,
				OrgID:       0,
				Name:        findData.Name,
				Des:         findData.Des,
				Permissions: vPermissions,
			}); err != nil {
				return errors.New("update user group, " + err.Error())
			}
		} else {
			if _, err := UserCore.CreateGroup(&UserCore.ArgsCreateGroup{
				OrgID:       0,
				Name:        v.Name,
				Des:         v.Des,
				Permissions: vPermissions,
			}); err != nil {
				return errors.New("create user group, " + err.Error())
			}
		}
	}

	//安装用户数据，如果存在最高权限则放弃
	type dataUserUserType struct {
		Name string
		//用户组的名称
		Groups   []string
		Username string
		Password string
	}
	type dataUsersType struct {
		Users []dataUserUserType
	}
	var userData dataUsersType
	if err := loadConfigFile("users.json", &userData); err != nil {
		return nil
	}
	for _, v := range userData.Users {
		findData, err := UserCore.GetUserByUsername(&UserCore.ArgsGetUserByUsername{
			OrgID:    0,
			Username: v.Username,
		})
		if err == nil && findData.ID > 0 {
			if Router2SystemConfig.Debug {
				if err := UserCore.UpdateUserInfoByID(&UserCore.ArgsUpdateUserInfoByID{
					ID:       findData.ID,
					OrgID:    -1,
					NewOrgID: -1,
					Name:     v.Name,
					Avatar:   0,
				},
				); err != nil {
					return errors.New("update user info, " + err.Error())
				}
				if err := UserCore.UpdateUserPasswordByID(&UserCore.ArgsUpdateUserPasswordByID{
					ID:       findData.ID,
					Password: v.Password,
				}); err != nil {
					return errors.New("update user username, " + err.Error())
				}
				if err := UserCore.UpdateUserUsernameByID(&UserCore.ArgsUpdateUserUsernameByID{
					ID: findData.ID, Username: v.Username,
				}); err != nil {
					return errors.New("update user username, " + err.Error())
				}
			}
		} else {
			var errCode string
			if findData, errCode, err = UserCore.CreateUser(&UserCore.ArgsCreateUser{
				OrgID:                0,
				Name:                 v.Name,
				Password:             v.Password,
				NationCode:           "",
				Phone:                "",
				AllowSkipPhoneVerify: false,
				AllowSkipWaitEmail:   false,
				Email:                "",
				Username:             v.Username,
				Avatar:               0,
				Status:               2,
				Parents:              nil,
				Groups:               nil,
				Infos:                nil,
				Logins:               nil,
				SortID:               0,
				Tags:                 nil,
			}); err != nil {
				return errors.New("create user, code: " + errCode + ", err: " + err.Error())
			}
		}
		//找到匹配的用户组，将批量更新
		for _, vGroupName := range v.Groups {
			var findGroupData UserCore.FieldsGroupType
			isFindGroup := false
			for _, vGroup := range allGroups {
				if vGroup.Name == vGroupName {
					isFindGroup = true
					findGroupData = vGroup
					break
				}
			}
			if !isFindGroup {
				continue
			}
			//更新用户组，管理员组
			if err := UserCore.UpdateUserGroupByID(&UserCore.ArgsUpdateUserGroupByID{
				ID:       findData.ID,
				OrgID:    0,
				GroupID:  findGroupData.ID,
				ExpireAt: time.Time{},
				IsRemove: false,
			}); err != nil {
				return errors.New("update user groups, " + err.Error())
			}
		}
		//如果启动debug
		if Router2SystemConfig.Debug {
			//检查用户是否被禁用？或其他非正常状态，自动修正
			if findData.Status != 2 {
				if err := UserCore.UpdateUserStatus(&UserCore.ArgsUpdateUserStatus{
					ID: findData.ID, Status: 2,
				}); err != nil {
					return errors.New("update user status, " + err.Error())
				}
			}
		}
	}

	//反馈
	return nil
}
