package BaseDBManager

import "time"

type FieldsSQL struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 来源
	// 如果存在值，尤其是带有FromCode时，应确保数据唯一性
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//来源系统
	// 例如: analysis
	FromSystem string `db:"from_system" json:"fromSystem" index:"true" field_list:"true"`
	//来源模块
	// 例如: index_sql
	FromModule string `db:"from_module" json:"fromModule" index:"true" field_list:"true"`
	//内部标识码
	// 可用于标记内部识别标识码，例如Index中的维度值，或一组维度值组合后的标识码
	FromCode string `db:"from_code" json:"fromCode" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 基础设置
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//定时器设置Carbon编码
	// 例如: 15s
	CarbonCode string `db:"carbon_code" json:"carbonCode" index:"true" field_list:"true"`
	//开始运行时通知中间件地址
	// 用于通知需发起该SQL，将SQL和来源信息传递给对应的中间件
	PostURL string `db:"post_url" json:"postURL" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 数据
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//SQL内容
	SQLData string `db:"sql_data" json:"sqlData"`
}
