package UserAddress

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetDefaultAddress 获取用户默认地址参数
type ArgsGetDefaultAddress struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetDefaultAddress 获取用户默认地址
func GetDefaultAddress(args *ArgsGetDefaultAddress) (data FieldsAddress, err error) {
	//获取默认地址数据
	var defaultData FieldsDefault
	err = Router2SystemConfig.MainDB.Get(&defaultData, "SELECT address_id FROM user_address_default WHERE user_id = $1", args.UserID)
	//找到则继续查询该地址
	if err == nil && defaultData.AddressID > 0 {
		data, err = GetIDTop(&ArgsGetIDTop{
			ID:     defaultData.AddressID,
			UserID: args.UserID,
		})
		if err == nil {
			//如果地址和实际地址不一致，说明地址发生过变更，修正默认值
			_ = SetDefault(&ArgsSetDefault{
				AddressID: data.ID,
				UserID:    args.UserID,
			})
			return
		}
		//不存在的地址，则继续
	}
	//抽取该用户任意一个地址反馈
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos FROM user_address WHERE user_id = $1 AND delete_at < to_timestamp(10000) LIMIT 1;", args.UserID)
	if err == nil {
		_ = SetDefault(&ArgsSetDefault{
			AddressID: data.ID,
			UserID:    args.UserID,
		})
	}
	return
}

// ArgsSetDefault 修改默认地址参数
type ArgsSetDefault struct {
	//地址ID
	AddressID int64 `db:"address_id" json:"addressID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// SetDefault 修改默认地址
func SetDefault(args *ArgsSetDefault) (err error) {
	//检查地址
	_, err = GetID(&ArgsGetID{
		ID:       args.AddressID,
		UserID:   args.UserID,
		IsRemove: false,
	})
	if err != nil {
		return
	}
	//获取记录
	var defaultData FieldsDefault
	err = Router2SystemConfig.MainDB.Get(&defaultData, "SELECT id, address_id FROM user_address_default WHERE user_id = $1", args.UserID)
	if err == nil {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_address_default SET update_at = NOW(), address_id = :address_id WHERE id = :id", map[string]interface{}{
			"id":         defaultData.ID,
			"address_id": args.AddressID,
		})
		return
	}
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_address_default (user_id, address_id) VALUES (:user_id, :address_id)", args)
	if err != nil {
		return
	}
	//更新数据
	OrgUserMod.PushUpdateUserData(0, args.UserID)
	//反馈
	return
}
