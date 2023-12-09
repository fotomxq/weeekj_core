package UserMessage

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsDeleteBySend 删除ID参数
type ArgsDeleteBySend struct {
	//IDs
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//发送人ID
	SendUserID int64 `db:"send_user_id" json:"sendUserID" check:"id" empty:"true"`
}

// DeleteBySend 删除ID
func DeleteBySend(args *ArgsDeleteBySend) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_message", "id = ANY(:ids) AND (:send_user_id < 1 OR send_user_id = :send_user_id)", args)
	if err != nil {
		return
	}
	for _, v := range args.IDs {
		deleteMessageCache(v)
	}
	return
}

// ArgsDeleteByReceive 接收人删除参数
type ArgsDeleteByReceive struct {
	//IDs
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//接收人
	ReceiveUserID int64 `db:"receive_user_id" json:"receiveUserID" check:"id" empty:"true"`
}

// DeleteByReceive 删除ID
func DeleteByReceive(args *ArgsDeleteByReceive) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_message SET receive_delete_at = NOW() WHERE id = ANY(:ids) AND (:receive_user_id < 1 OR receive_user_id = :receive_user_id) AND status = 2 AND receive_delete_at < to_timestamp(1000000)", args)
	if err != nil {
		return
	}
	for _, v := range args.IDs {
		deleteMessageCache(v)
	}
	return
}

// ArgsDeleteByID 接收人删除参数
type ArgsDeleteByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteByID 删除ID
func DeleteByID(args *ArgsDeleteByID) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_message SET delete_at = NOW(), receive_delete_at = NOW() WHERE id = :id", args)
	if err != nil {
		return
	}
	deleteMessageCache(args.ID)
	return
}
