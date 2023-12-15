package ServiceInfoExchange

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteTake 删除报名参数
type ArgsDeleteTake struct {
	//信息ID
	ID int64 `db:"id" json:"id" check:"id"`
	//信息所属用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//报名人用户ID
	TakeUserID int64 `json:"takeUserID" check:"id"`
}

// DeleteTake 删除报名
func DeleteTake(args *ArgsDeleteTake) (err error) {
	infoData := getInfoByID(args.ID)
	if infoData.ID < 1 || !CoreFilter.EqID2(args.UserID, infoData.UserID) {
		err = errors.New("no data")
		return
	}
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "service_info_exchange_take", "info_id = :info_id AND user_id = :user_id", map[string]interface{}{
		"info_id": infoData.ID,
		"user_id": args.TakeUserID,
	})
	return
}
