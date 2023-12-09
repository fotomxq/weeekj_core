package UserAddress

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsAddress 用户地址信息
type FieldsAddress struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//上级
	// 如果存在上级，则说明为历史数据
	ParentID int64 `db:"parent_id" json:"parentID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//地址昵称
	NiceName string `db:"nice_name" json:"niceName"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province"`
	//所属城市
	City int `db:"city" json:"city"`
	//街道详细信息
	Address string `db:"address" json:"address"`
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//联系人姓名
	Name string `db:"name" json:"name"`
	//联系人国家代码
	NationCode string `db:"nation_code" json:"nationCode"`
	//联系人手机号
	Phone string `db:"phone" json:"phone"`
	//联系人邮箱
	Email string `db:"email" json:"email"`
	//其他联系方式
	Infos CoreSQLConfig.FieldsInfosType `db:"infos" json:"infos"`
}

// FieldsDefault 用户默认表
type FieldsDefault struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//地址ID
	AddressID int64 `db:"address_id" json:"addressID"`
}
