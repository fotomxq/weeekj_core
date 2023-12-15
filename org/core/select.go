package OrgCoreCore

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsSetSelect 设置用户选择某个组织参数
type ArgsSetSelect struct {
	//用户ID
	UserID int64 `json:"userID"`
	//组织ID
	OrgID int64 `json:"orgID"`
}

// SetSelect 设置用户选择某个组织
// 内部将检查权限及能否选择该组织
func SetSelect(args *ArgsSetSelect) (bindData FieldsBind, permissions []string, err error) {
	//获取用户在该组织的绑定关系
	bindData, err = GetBindByUserAndOrg(&ArgsGetBindByUserAndOrg{
		UserID: args.UserID,
		OrgID:  args.OrgID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get bind by create info and org, ", err))
		return
	}
	permissions = bindData.Manager
	//获取分组数据
	var groupList []FieldsGroup
	groupList, err = getGroups(bindData.GroupIDs)
	if err == nil {
		for _, v := range groupList {
			permissions = margeBindManager(permissions, v.Manager)
		}
	} else {
		//不用记录错误
		err = nil
	}
	//必须具备至少是成员的权限
	isFind := false
	for _, v := range permissions {
		if v == "member" || v == "all" {
			isFind = true
			break
		}
	}
	if !isFind {
		err = errors.New("no permission")
		return
	}
	//如果存在则构建绑定关系
	var data FieldsSelectOrg
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_select WHERE user_id = $1;", args.UserID)
	if err == nil {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_select SET last_at = NOW(), org_id = :org_id, bind_id = :bind_id WHERE id = :id;", map[string]interface{}{
			"id":      data.ID,
			"org_id":  bindData.OrgID,
			"bind_id": bindData.ID,
		})
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_core_select (user_id, org_id, bind_id) VALUES (:user_id, :org_id, :bind_id)", map[string]interface{}{
			"user_id": args.UserID,
			"org_id":  bindData.OrgID,
			"bind_id": bindData.ID,
		})
	}
	if err == nil {
		//更新最后1次登陆时间
		err = updateBindLastTime(&argsUpdateBindLastTime{
			ID: bindData.ID,
		})
		if err != nil {
			return
		}
	}
	return
}

// ArgsGetSelect 获取用户当前选择的组织参数
type ArgsGetSelect struct {
	//用户ID
	UserID int64 `json:"userID"`
}

// GetSelect 获取用户当前选择的组织
func GetSelect(args *ArgsGetSelect) (data FieldsSelectOrg, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, last_at, user_id, org_id, bind_id FROM org_core_select WHERE user_id = $1 ORDER BY id DESC LIMIT 1", args.UserID)
	return
}

// ArgsGetSelectAndData 获取用户当前选择的组织，同时反馈其他数据参数
type ArgsGetSelectAndData struct {
	//用户ID
	UserID int64 `json:"userID"`
}

// GetSelectAndData 获取用户当前选择的组织，同时反馈其他数据参数
func GetSelectAndData(args *ArgsGetSelectAndData) (orgData FieldsOrg, bindData FieldsBind, permissions []string, err error) {
	var selectData FieldsSelectOrg
	selectData, err = GetSelect(&ArgsGetSelect{
		UserID: args.UserID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get select, ", err))
		return
	}
	orgData, err = GetOrg(&ArgsGetOrg{
		ID: selectData.OrgID,
	})
	if err != nil || CoreSQL.CheckTimeHaveData(orgData.DeleteAt) {
		err = errors.New(fmt.Sprint("get org, ", err))
		return
	}
	bindData, err = GetBind(&ArgsGetBind{
		ID:     selectData.BindID,
		OrgID:  -1,
		UserID: -1,
	})
	if err != nil || CoreSQL.CheckTimeHaveData(bindData.DeleteAt) {
		err = errors.New(fmt.Sprint("get bind, ", err))
		return
	}
	//获取权限
	permissions = GetPermissionByBindID(bindData.ID)
	//反馈
	return
}

// 组合权限工具
func margeBindManager(permissionA []string, permissionB []string) []string {
	isFind := false
	for _, v := range permissionB {
		for _, v2 := range permissionA {
			if v == v2 {
				isFind = true
				break
			}
		}
		if !isFind {
			permissionA = append(permissionA, v)
		}
	}
	return permissionA
}
