package UserTicketSend

import (
	"errors"
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetSendListByUser 获取用户可领取的票据列表参数
type ArgsGetSendListByUser struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// DataGetSendListByUser 获取用户可领取的票据列表数据
type DataGetSendListByUser struct {
	//发放ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//发放的票据配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//是否必须是会员配置ID
	NeedUserSubConfigID int64 `db:"need_user_sub_config_id" json:"needUserSubConfigID"`
	//是否自动发放，如果不是，则需绑定广告
	NeedAuto bool `db:"need_auto" json:"needAuto"`
	//每个用户发放几张
	PerCount int64 `db:"per_count" json:"perCount"`
}

// GetSendListByUser 获取用户可领取的票据列表
func GetSendListByUser(args *ArgsGetSendListByUser) (dataList []DataGetSendListByUser, err error) {
	var sendList []FieldsSend
	err = Router2SystemConfig.MainDB.Select(&sendList, "SELECT id, org_id, config_id, need_user_sub_config_id, need_auto, per_count FROM user_ticket_send WHERE ($1 < 0 OR org_id = $1) AND finish_at < to_timestamp(1000000) AND need_auto = false LIMIT 1000", args.OrgID)
	if err != nil || len(sendList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range sendList {
		if checkSendTakeByUser(args.UserID, v.ID) {
			continue
		}
		dataList = append(dataList, DataGetSendListByUser{
			ID:                  v.ID,
			OrgID:               v.OrgID,
			ConfigID:            v.ConfigID,
			NeedUserSubConfigID: v.NeedUserSubConfigID,
			NeedAuto:            v.NeedAuto,
			PerCount:            v.PerCount,
		})
	}
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsTakeSend 领取优惠券参数
type ArgsTakeSend struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//要领取的SendID
	SendID int64 `db:"send_id" json:"sendID" check:"id"`
}

// TakeSend 领取优惠券
func TakeSend(args *ArgsTakeSend) (errCode string, err error) {
	//获取配置
	var sendData FieldsSend
	err = Router2SystemConfig.MainDB.Get(&sendData, "SELECT id, create_at, finish_at, org_id, send_count, need_user_sub_config_id, need_auto, config_id, per_count FROM user_ticket_send WHERE id = $1 AND need_auto = false AND finish_at < to_timestamp(1000000) AND ($2 < 1 OR org_id = $2)", args.SendID, args.OrgID)
	if err != nil || sendData.ID < 1 {
		errCode = "err_no_data"
		err = errors.New(fmt.Sprint("send not find by id: ", args.SendID, ", org id: ", args.OrgID))
		return
	}
	//检查用户是否领取？
	if checkSendTakeByUser(args.UserID, sendData.ID) {
		errCode = "err_user_have_take"
		err = errors.New(fmt.Sprint("user have log, user id: ", args.UserID, ", send id: ", args.SendID))
		return
	}
	//开始赠送
	if !runSendUser(args.UserID, &sendData) {
		errCode = "err_update"
		err = errors.New(fmt.Sprint("user no sub or other error, user id: ", args.UserID, ", send id: ", sendData.ID))
		return
	}
	//反馈
	return
}

// checkSendTakeByUser 检查用户是否已经领取
func checkSendTakeByUser(userID int64, sendID int64) (isTake bool) {
	var logID int64
	err := Router2SystemConfig.MainDB.Get(&logID, "SELECT id FROM user_ticket_send_log WHERE user_id = $1 AND send_id = $2 LIMIT 1", userID, sendID)
	return err == nil && logID > 0
}
