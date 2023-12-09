package UserCore

import (
	"encoding/json"
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 检查用户是否可以成为上级
// 递归检查上级关系，如果形成悖论，则反馈失败
func checkUserParentID(id int64, parentID int64, systemMark string, checkParentIDs []int64) error {
	checkParentIDs = append(checkParentIDs, id)
	if id == parentID {
		return errors.New("parent id is cycle")
	}
	parentData, err := GetUserByID(&ArgsGetUserByID{
		ID: parentID,
	})
	if err != nil {
		//上级不存在
		return errors.New("parent is not exist")
	}
	//上级存在，则检查其上级是否在序列内？
	for _, v := range checkParentIDs {
		for _, v2 := range parentData.Parents {
			if v == v2.ParentID && systemMark == v2.System {
				return errors.New("parent id is cycle")
			}
		}
	}
	//不存在上级，则跳出
	if len(parentData.Parents) < 1 {
		return nil
	}
	//不存在则继续
	for _, v := range parentData.Parents {
		if err := checkUserParentID(parentData.ID, v.ParentID, systemMark, checkParentIDs); err != nil {
			return err
		}
	}
	return nil
}

// 获取用户组使用数量
func getUserCountByGroup(groupID int64, orgID int64) (count int64) {
	type argType struct {
		GroupID int64 `json:"group_id"`
	}
	argsJSON, err := json.Marshal([]argType{
		{
			GroupID: groupID,
		},
	})
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "user_core", "id", "org_id = $1 AND groups @> $2", orgID, string(argsJSON))
	if err != nil {
		return
	}
	return
}

// 检查密码
func checkPassword(password string) bool {
	return CoreFilter.CheckPassword(password)
}

// 获取密码摘要
func getPasswordSha(password string) (string, error) {
	return CoreFilter.GetSha1ByString(password)
}

// 检查status
func checkUserStatus(status int) error {
	switch status {
	case 0:
	case 1:
	case 2:
	default:
		return errors.New("status is error")
	}
	return nil
}
