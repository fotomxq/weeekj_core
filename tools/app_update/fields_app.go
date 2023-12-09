package ToolsAppUpdate

import (
	"github.com/lib/pq"
	"time"
)

//FieldsApp APP
type FieldsApp struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name" json:"name"`
	//升级内容
	Des string `db:"des" json:"des"`
	//描述附带文件
	DesFiles  pq.Int64Array `db:"des_files" json:"desFiles"`
	//应用标识码
	AppMark string `db:"app_mark" json:"appMark"`
	//总下载次数
	Count int64 `db:"count" json:"count"`
}