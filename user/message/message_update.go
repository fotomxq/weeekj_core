package UserMessage

import (
	"errors"
	"fmt"
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsUpdateByID 更新消息内容参数
// 只有发送人在草稿状态下可以编辑消息
type ArgsUpdateByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//发送人ID
	// 用于验证
	SendUserID int64 `db:"send_user_id" json:"sendUserID" check:"id" empty:"true"`
	//接收人
	ReceiveUserID int64 `db:"receive_user_id" json:"receiveUserID" check:"id"`
	//预计发送时间
	WaitSendAt time.Time `db:"wait_send_at" json:"waitSendAt" check:"isoTime" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"des" min:"1" max:"300"`
	//内容
	Content string `db:"content" json:"content" check:"des" min:"1" max:"1500"`
	//附件文件列
	Files pq.Int64Array `db:"files" json:"files" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateByID 更新消息内容
// 必须是草稿状态
func UpdateByID(args *ArgsUpdateByID) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_message SET update_at = NOW(), wait_send_at = :wait_send_at, title = :title, content = :content, files = :files, params = :params WHERE id = :id AND (:send_user_id < 1 OR send_user_id = :send_user_id) AND (:receive_user_id < 1 OR receive_user_id = :receive_user_id) AND status = 0", args)
	if err != nil {
		return
	}
	deleteMessageCache(args.ID)
	return
}

// ArgsUpdatePost 推送提交参数
type ArgsUpdatePost struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//发送人ID
	// 用于验证
	SendUserID int64 `db:"send_user_id" json:"sendUserID" check:"id" empty:"true"`
}

// UpdatePost 推送提交
func UpdatePost(args *ArgsUpdatePost) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_message SET status = 1 WHERE id = :id AND (:send_user_id < 1 OR send_user_id = :send_user_id) AND status = 0", args)
	if err != nil {
		return
	}
	deleteMessageCache(args.ID)
	//当系统不需要人工审核时，进行自动审核处理
	userMessageAudit, _ := BaseConfig.GetDataBool("UserMessageAudit")
	if err != nil {
		userMessageAudit = true
	}
	if userMessageAudit {
		pushNatsAutoAudit(args.ID)
	}
	return
}

// ArgsUpdateAudit 审核消息参数
type ArgsUpdateAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// UpdateAudit 审核消息
func UpdateAudit(args *ArgsUpdateAudit) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_message SET status = 2 WHERE id = :id AND status = 1", args)
	if err != nil {
		return
	}
	deleteMessageCache(args.ID)
	return
}

// ArgsUpdateReceiveRead 已经阅读参数
type ArgsUpdateReceiveRead struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//接收人
	ReceiveUserID int64 `db:"receive_user_id" json:"receiveUserID" check:"id" empty:"true"`
}

// UpdateReceiveRead 已经阅读
func UpdateReceiveRead(args *ArgsUpdateReceiveRead) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_message SET receive_read_at = NOW() WHERE id = :id AND (:receive_user_id < 1 OR receive_user_id = :receive_user_id) AND status = 2 AND receive_read_at < to_timestamp(1000000)", args)
	if err != nil {
		return
	}
	deleteMessageCache(args.ID)
	if args.ReceiveUserID > 0 {
		var countUnRead int64
		_ = Router2SystemConfig.MainDB.Get(&countUnRead, "SELECT COUNT(id) FROM user_message WHERE receive_user_id = $1 AND receive_delete_at < to_timestamp(1000000) AND receive_read_at < to_timestamp(1000000)", args.ReceiveUserID)
		AnalysisAny2.AppendData("re", "user_message_receive_unread_count", time.Time{}, 0, args.ReceiveUserID, args.ReceiveUserID, 0, 0, countUnRead)
		var countRead int64
		_ = Router2SystemConfig.MainDB.Get(&countRead, "SELECT COUNT(id) FROM user_message WHERE receive_user_id = $1 AND receive_read_at > to_timestamp(1000000)", args.ReceiveUserID)
		AnalysisAny2.AppendData("re", "user_message_receive_read_count", time.Time{}, 0, args.ReceiveUserID, args.ReceiveUserID, 0, 0, countRead)
	}
	return
}

// ArgsUpdateReceiveReads 批量设置已读参数
type ArgsUpdateReceiveReads struct {
	//IDs
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//接收人
	ReceiveUserID int64 `db:"receive_user_id" json:"receiveUserID" check:"id" empty:"true"`
}

// UpdateReceiveReads 批量设置已读
func UpdateReceiveReads(args *ArgsUpdateReceiveReads) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_message SET receive_read_at = NOW() WHERE id = ANY(:ids) AND (:receive_user_id < 1 OR receive_user_id = :receive_user_id) AND status = 2 AND receive_read_at < to_timestamp(1000000)", args)
	if err != nil {
		err = errors.New(fmt.Sprint("update failed, ", err, ", ids: ", args.IDs, ", receive_user_id: ", args.ReceiveUserID))
		return
	}
	for _, v := range args.IDs {
		deleteMessageCache(v)
	}
	return
}
