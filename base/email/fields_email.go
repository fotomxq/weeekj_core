package BaseEmail

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsEmailType 等待发送的数据
type FieldsEmailType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//使用的发送渠道ID
	ServerID int64 `db:"server_id" json:"serverID"`
	//创建来源和创建来源ID
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//预计发送时间
	SendAt time.Time `db:"send_at" json:"sendAt"`
	//是否已经发送？
	IsSuccess bool `db:"is_success" json:"isSuccess"`
	//是否发送失败
	IsFailed bool `db:"is_failed" json:"isFailed"`
	//失败原因
	FailMessage string `db:"fail_message" json:"failMessage"`
	//标题
	Title string `db:"title" json:"title"`
	//内容
	Content string `db:"content" json:"content"`
	//文本类型
	// text / html
	ContentType string `db:"content_type" json:"contentType"`
	//目标人
	ToEmail string `db:"to_email" json:"toEmail"`
}
