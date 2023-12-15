package UserCore

import (
	"errors"
	"fmt"
	BaseQiniu "github.com/fotomxq/weeekj_core/v5/base/qiniu"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
)

// ArgsGetUserData 获取用户的整合信息参数
type ArgsGetUserData struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
}

// Deprecated: 准备废弃
// GetUserData 获取用户的整合信息
func GetUserData(args *ArgsGetUserData) (userData DataUserDataType, err error) {
	//获取用户基本信息
	userData.Info, err = GetUserByID(&ArgsGetUserByID{
		ID:    args.UserID,
		OrgID: -1,
	})
	if err != nil {
		err = errors.New("find user info, " + err.Error())
		return
	}
	//检查用户状态？
	if userData.Info.Status != 2 {
		err = errors.New("status no public")
		return
	}
	//获取文件列
	var fileIDs []int64
	if userData.Info.Avatar > 0 {
		fileIDs = append(fileIDs, userData.Info.Avatar)
	}
	if len(fileIDs) > 0 {
		userData.FileList, err = BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
			ClaimIDList: fileIDs,
			UserID:      0,
			OrgID:       0,
			IsPublic:    true,
		})
		if err != nil {
			err = nil
		}
	}
	//获取组织可用的所有权限集
	var orgPermission []FieldsPermissionType
	orgPermission, err = GetAllPermission(&ArgsGetAllPermission{
		AllowOrg: userData.Info.OrgID > 0,
	})
	//遍历用户组，并根据过期时间插入数据
	nowTime := CoreFilter.GetNowTime().Unix()
	for _, v := range userData.Info.Groups {
		if v.ExpireAt.Unix() < 1 || v.ExpireAt.Unix() >= nowTime {
			groupInfo, err := GetGroup(&ArgsGetGroup{
				ID: v.GroupID,
			})
			if err != nil {
				//找不到group，可能被删除了，将不会被加入到用户组中
				continue
			}
			//匹配组织
			// 禁止越权访问数据
			// 如果该组属于组织，则判断该设计；如果属于平台，则不判断。
			if groupInfo.OrgID > 0 && groupInfo.OrgID != userData.Info.OrgID {
				continue
			}
			//写入数据集合
			userData.Groups = append(userData.Groups, groupInfo)
			//查询该用户组的权限，检查并插入info
			for _, v2 := range groupInfo.Permissions {
				//检查组织是否可以授权，如果不能则跳过
				if groupInfo.OrgID > 0 && userData.Info.OrgID > 0 {
					allowOrg := false
					for _, v3 := range orgPermission {
						if v3.AllowOrg {
							allowOrg = true
							break
						}
					}
					if !allowOrg {
						continue
					}
				}
				//写入列队
				userData.Permissions = appendUserPermission(userData.Permissions, v2)
			}
		}
	}
	//反馈
	return
}

// ArgsCheckUserPassword 验证用户名和密码参数
type ArgsCheckUserPassword struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户信息
	UserInfo *FieldsUserType
	//要验证的密码
	Password string
}

// CheckUserPassword 验证用户名和密码
func CheckUserPassword(args *ArgsCheckUserPassword) (err error) {
	if args.OrgID > -1 {
		if !CoreFilter.EqID2(args.OrgID, args.UserInfo.OrgID) {
			err = errors.New(fmt.Sprint("not this org, find id: ", args.OrgID, ", user org id: ", args.UserInfo.OrgID))
			return
		}
	}
	var passwordSha1 string
	passwordSha1, err = getPasswordSha(args.Password)
	if err != nil {
		return err
	}
	if args.UserInfo.Password != passwordSha1 {
		err = errors.New("password hash is error")
		return
	}
	return
}

// ArgsFilterUserData 脱敏处理用户结构体参数
type ArgsFilterUserData struct {
	//用户信息
	UserData DataUserDataType
}

// Deprecated: 准备废弃
// FilterUserData 脱敏处理用户结构体参数
func FilterUserData(args *ArgsFilterUserData) DataUserDataType {
	info := FilterUserInfo(&ArgsFilterUserInfo{
		UserInfo: args.UserData.Info,
	})
	if info.Phone != "" && len(info.Phone) == 11 {
		info.Phone = string(info.Phone[0]) + string(info.Phone[1]) + string(info.Phone[2]) + "***"
	}
	return DataUserDataType{
		Info:        info,
		Groups:      args.UserData.Groups,
		Permissions: args.UserData.Permissions,
		FileList:    args.UserData.FileList,
	}
}

// ArgsFilterUserInfo 用户脱敏参数
type ArgsFilterUserInfo struct {
	//用户信息
	UserInfo FieldsUserType
}

// FilterUserInfo 用户脱敏
func FilterUserInfo(args *ArgsFilterUserInfo) FieldsUserType {
	return FieldsUserType{
		ID:         args.UserInfo.ID,
		CreateAt:   args.UserInfo.CreateAt,
		UpdateAt:   args.UserInfo.UpdateAt,
		DeleteAt:   args.UserInfo.DeleteAt,
		Status:     args.UserInfo.Status,
		OrgID:      args.UserInfo.OrgID,
		Name:       args.UserInfo.Name,
		Password:   "",
		NationCode: args.UserInfo.NationCode,
		Phone:      args.UserInfo.Phone,
		Email:      args.UserInfo.Email,
		Username:   args.UserInfo.Username,
		Avatar:     args.UserInfo.Avatar,
		Parents:    args.UserInfo.Parents,
		Groups:     args.UserInfo.Groups,
		Infos:      args.UserInfo.Infos,
		Logins:     []FieldsUserLoginType{},
	}
}

// GetFromByUser 组合from来源结构
func GetFromByUser(data *DataUserDataType) CoreSQLFrom.FieldsFrom {
	return CoreSQLFrom.FieldsFrom{
		System: "user",
		ID:     data.Info.ID,
		Mark:   "",
		Name:   data.Info.Name,
	}
}

// appendUserPermission 汇集用户权限，如果存在则跳过
func appendUserPermission(permissions []string, mark string) []string {
	isFind := false
	for _, v := range permissions {
		if v == mark {
			isFind = true
			break
		}
	}
	if !isFind {
		permissions = append(permissions, mark)
	}
	return permissions
}
