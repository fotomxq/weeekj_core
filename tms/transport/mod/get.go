package TMSTransportMod

import Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"

// ArgsGetTransport 获取配送信息参数
type ArgsGetTransport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetTransport 获取配送信息
func GetTransport(args *ArgsGetTransport) (data FieldsTransport, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, finish_at, org_id, bind_id, info_id, user_id, sn, sn_day, status, from_address, to_address, order_id, goods, weight, length, width, currency, price, pay_finish_at, pay_id, pay_ids, task_at, params FROM tms_transport WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR info_id = $3) AND ($4 < 1 OR user_id = $4)", args.ID, args.OrgID, args.InfoID, args.UserID)
	return
}
