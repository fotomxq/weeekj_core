package OrgUser

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserAddress "github.com/fotomxq/weeekj_core/v5/user/address"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

// ArgsGetUserDataList 获取用户数据列表参数
type ArgsGetUserDataList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//搜索电话
	SearchPhone string `json:"searchPhone" check:"search" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetUserDataList 获取用户数据列表
func GetUserDataList(args *ArgsGetUserDataList) (dataList []FieldsOrgUser, dataCount int64, err error) {
	where := "org_id = $1 AND ($2 = '' OR (phone = $2 OR address_list -> 0 -> 'phone' ? $2)) AND ($3 = '' OR (address_list::text ILIKE '%' || $3 || '%' OR name ILIKE '%' || $3 || '%' OR phone ILIKE '%' || $3 || '%'))"
	tableName := "org_user_data"
	var rawList []FieldsOrgUser
	dataCount, err = CoreSQL.GetListPageAndCountArgs(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id, org_id, user_id FROM "+tableName+" WHERE "+where,
		where,
		&args.Pages,
		[]string{"id", "create_at", "update_at"},
		args.OrgID,
		args.SearchPhone,
		args.Search,
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getUserData(v.OrgID, v.UserID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetUserData 获取用户数据集合参数
type ArgsGetUserData struct {
	//组织ID
	OrgID int64 `json:"orgID" check:"id"`
	//搜索电话
	SearchPhone string `json:"searchPhone" check:"search" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetUserData 获取用户数据集合
func GetUserData(args *ArgsGetUserData) (dataList []FieldsOrgUser, err error) {
	//从数据集合中搜索
	err = Router2SystemConfig.MainDB.Select(&dataList, fmt.Sprint("SELECT id, create_at, update_at, org_id, user_id, name, phone, address_list, user_integral, user_subs, user_tickets, deposit_data, last_order, params FROM org_user_data WHERE org_id = $1 AND ($2 = '' OR (phone = $2 OR address_list @> '{\"phone\": \"", args.SearchPhone, "\"}')) AND ($3 = '' OR (name ILIKE '%' || $3 || '%' OR address_list::text ILIKE '%' || $3 || '%')) ORDER BY update_at DESC LIMIT 10"), args.OrgID, args.SearchPhone, args.Search)
	if err == nil && len(dataList) > 0 {
		return
	}
	//找不到则更新数据
	//在用户中查询该用户
	searchUser := args.SearchPhone
	if searchUser == "" {
		searchUser = args.Search
	}
	var userList []UserCore.FieldsUserType
	userList, _, err = UserCore.GetUserList(&UserCore.ArgsGetUserList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		Status:       2,
		OrgID:        args.OrgID,
		ParentSystem: "",
		ParentID:     0,
		SortID:       -1,
		Tags:         []int64{},
		IsRemove:     false,
		Search:       searchUser,
	})
	if err == nil && len(userList) > 0 {
		for _, v := range userList {
			var vData FieldsOrgUser
			err = updateByUserID(args.OrgID, v.ID)
			if err != nil {
				err = nil
				continue
			}
			vData = getUserData(args.OrgID, v.ID)
			if vData.ID < 1 {
				continue
			}
			dataList = append(dataList, vData)
		}
		if len(dataList) > 0 {
			return
		}
	}
	//从地址中查询
	var addressList []UserAddress.FieldsAddress
	addressList, _, err = UserAddress.GetList(&UserAddress.ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		ParentID:    0,
		UserID:      0,
		Country:     0,
		Province:    0,
		City:        0,
		IsRemove:    false,
		SearchPhone: args.SearchPhone,
		Search:      args.Search,
	})
	if err == nil && len(addressList) < 1 {
		for _, v := range addressList {
			var vData FieldsOrgUser
			err = updateByUserID(args.OrgID, v.UserID)
			if err != nil {
				err = nil
				continue
			}
			vData = getUserData(args.OrgID, v.UserID)
			if vData.ID < 1 {
				continue
			}
			dataList = append(dataList, vData)
		}
		if len(dataList) > 0 {
			return
		}
	}
	//反馈数据
	if len(dataList) < 1 {
		err = errors.New("data is empty")
		return
	}
	return
}

// ArgsGetUserDataByUserID 通过用户ID获取聚合数据参数
type ArgsGetUserDataByUserID struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetUserDataByUserID 通过用户ID获取聚合数据
func GetUserDataByUserID(args *ArgsGetUserDataByUserID) (data FieldsOrgUser, err error) {
	data = getUserData(args.OrgID, args.UserID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		return
	}
	return
}

// 获取用户数据
func getUserData(orgID int64, userID int64) (data FieldsOrgUser) {
	cacheMark := getUserCacheMark(orgID, userID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, org_id, user_id, name, phone, address_list, user_integral, user_subs, user_tickets, deposit_data, last_order, params FROM org_user_data WHERE org_id = $1 AND user_id = $2 LIMIT 1", orgID, userID)
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 43200)
	return
}

// ArgsCheckUserDataByOrg 检查组织下是否具备对应用户聚合数据参数
type ArgsCheckUserDataByOrg struct {
	//组织ID
	OrgID int64 `json:"orgID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// CheckUserDataByOrg 检查组织下是否具备对应用户聚合数据
func CheckUserDataByOrg(args *ArgsCheckUserDataByOrg) (b bool) {
	var data FieldsOrgUser
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_user_data WHERE org_id = $1 AND user_id = $2 LIMIT 1", args.OrgID, args.UserID)
	if err == nil && data.ID > 0 {
		b = true
		return
	}
	return
}

func getUserDataIDByUserID(orgID int64, userID int64) (id int64) {
	var data FieldsOrgUser
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_user_data WHERE org_id = $1 AND user_id = $2 LIMIT 1", orgID, userID)
	if err == nil && data.ID > 0 {
		id = data.ID
		return
	}
	return
}
