package UserRecordCore

import (
	"time"
)

//FieldsRecordType 短信模版和配置信息结构
type FieldsRecordType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//用户昵称
	UserName string `db:"username" json:"userName"`
	//操作内容标识码
	// 可用于其他语言处理
	ContentMark string `db:"content_mark" json:"contentMark"`
	//操作内容概述
	Content string `db:"content" json:"content"`
}
