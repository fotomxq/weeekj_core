package DataLakeSource

import "time"

// FieldsTable 基础表信息
// 注意不要将多个来源混合到一个表中，应拆分表；可以将一个来源根据需求，进行拆分，但不推荐
type FieldsTable struct {
	//ID
	ID int64 `db:"id" json:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0"`
	//表名称
	TableName string `db:"table_name" json:"tableName" index:"true" field_search:"true"`
	//表描述
	TableDesc string `db:"table_desc" json:"tableDesc" field_search:"true"`
	//提示名称
	TipName string `db:"tip_name" json:"tipName" field_search:"true"`
	//数据唯一渠道名称
	// 如果是多处来源，应拆分表
	ChannelName string `db:"channel_name" json:"channelName" field_search:"true"`
	//数据唯一渠道提示名称
	ChannelTipName string `db:"channel_tip_name" json:"channelTipName" field_search:"true"`
}
