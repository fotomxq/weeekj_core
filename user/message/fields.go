package UserMessage

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

type FieldsMessage struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	// 发送人删除状态
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//预计发送时间
	WaitSendAt time.Time `db:"wait_send_at" json:"waitSendAt"`
	//发送状态
	// 发送完成后，发送人无法删除，接收人可以标记已读或删除，具体其他字段完成该操作
	// 0 草稿; 1 等待审核; 2 发送成功
	Status int `db:"status" json:"status"`
	//发送人
	SendUserID int64 `db:"send_user_id" json:"sendUserID"`
	//接收人
	ReceiveUserID int64 `db:"receive_user_id" json:"receiveUserID"`
	//接收人阅读时间
	ReceiveReadAt time.Time `db:"receive_read_at" json:"receiveReadAt"`
	//接收人删除状态
	ReceiveDeleteAt time.Time `db:"receive_delete_at" json:"receiveDeleteAt"`
	//标题
	Title string `db:"title" json:"title"`
	//内容
	Content string `db:"content" json:"content"`
	//附件文件列
	Files pq.Int64Array `db:"files" json:"files"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
