package ToolsCommunication

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

type FieldsRoom struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//到期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务字符串; 3 第三方agora服务uint32; 4 第三方agora服务字符串trc; 5 第三方agora服务uint32 rtc
	ConnectType int `db:"connect_type" json:"connectType"`
	//通讯类型
	DataType int `db:"data_type" json:"dataType"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//房间名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//是否公开房间？
	// 私有化房间只允许特定链接链接，否则可以通过公共列表查询到
	IsPublic bool `db:"is_public" json:"isPublic"`
	//房间链接密码
	Password string `db:"password" json:"password"`
	//最大人数
	MaxCount int `db:"max_count" json:"maxCount"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
