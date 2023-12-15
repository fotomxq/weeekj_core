package OrgUser

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsSearchUser 搜索指定的用户，或锁定ID参数
type ArgsSearchUser struct {
	//组织ID
	OrgID int64 `json:"orgID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//搜索电话
	SearchPhone string `json:"searchPhone" check:"search" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// SearchUser 搜索指定的用户，或锁定ID
func SearchUser(args *ArgsSearchUser) (dataList []FieldsOrgUser) {
	//如果存在ID，则找到直接锁定
	if args.OrgID > 0 && args.UserID > 0 {
		data := getUserData(args.OrgID, args.UserID)
		if data.ID > 0 {
			dataList = append(dataList, data)
		}
	}
	//搜索符合条件的数据
	where := "org_id = $1 AND ($2 = '' OR (phone = $2 OR address_list -> 0 -> 'phone' ? $2 OR address_list -> 1 -> 'phone' ? $2 OR address_list -> 2 -> 'phone' ? $2)) AND ($3 = '' OR (address_list::text ILIKE '%' || $3 || '%' OR name ILIKE '%' || $3 || '%' OR phone ILIKE '%' || $3 || '%'))"
	tableName := "org_user_data"
	var rawList []FieldsOrgUser
	_, err := CoreSQL.GetListPageAndCountArgs(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id, org_id, user_id FROM "+tableName+" WHERE "+where,
		where,
		&CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  20,
			Sort: "update_at",
			Desc: true,
		},
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
		if len(vData.AddressList) < 1 {
			if vData.LastOrder.AddressTo.Address != "" {
				vData.AddressList = append(vData.AddressList, FieldsAddress{
					ID:         0,
					UpdateAt:   time.Time{},
					Country:    vData.LastOrder.AddressTo.Country,
					Province:   vData.LastOrder.AddressTo.Province,
					City:       vData.LastOrder.AddressTo.City,
					Address:    vData.LastOrder.AddressTo.Address,
					MapType:    vData.LastOrder.AddressTo.MapType,
					Longitude:  vData.LastOrder.AddressTo.Longitude,
					Latitude:   vData.LastOrder.AddressTo.Latitude,
					Name:       vData.LastOrder.AddressTo.Name,
					NationCode: vData.LastOrder.AddressTo.NationCode,
					Phone:      vData.LastOrder.AddressTo.Phone,
				})
			}
		}
		dataList = append(dataList, vData)
	}
	return
}
