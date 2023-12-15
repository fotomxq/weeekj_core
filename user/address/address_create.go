package UserAddress

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCreate 创建新的地址参数
type ArgsCreate struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//地址昵称
	NiceName string `db:"nice_name" json:"niceName" check:"name" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province" check:"province"`
	//所属城市
	City int `db:"city" json:"city" check:"city"`
	//街道详细信息
	Address string `db:"address" json:"address" check:"address"`
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType" check:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude" check:"gps"`
	Latitude  float64 `db:"latitude" json:"latitude" check:"gps"`
	//联系人姓名
	Name string `db:"name" json:"name" check:"name"`
	//联系人国家代码
	NationCode string `db:"nation_code" json:"nationCode" check:"nationCode"`
	//联系人手机号
	Phone string `db:"phone" json:"phone" check:"phone"`
	//联系人邮箱
	Email string `db:"email" json:"email" check:"email" empty:"true"`
	//其他联系方式
	Infos CoreSQLConfig.FieldsInfosType `db:"infos" json:"infos"`
}

// Create 创建新的地址
func Create(args *ArgsCreate) (data FieldsAddress, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_address", "INSERT INTO user_address (delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos) VALUES (to_timestamp(0), 0, :user_id, :nice_name, :country, :province, :city, :address, :map_type, :longitude, :latitude, :name, :nation_code, :phone, :email, :infos)", args, &data)
	if err != nil {
		return
	}
	//更新组织用户数据包
	OrgUserMod.PushUpdateUserData(0, args.UserID)
	//反馈
	return
}
