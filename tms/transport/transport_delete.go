package TMSTransport

import (
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteTransport 删除配送单参数
type ArgsDeleteTransport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//描述信息
	Des string `json:"des" check:"des" min:"1" max:"1000"`
}

// DeleteTransport 删除配送单
func DeleteTransport(args *ArgsDeleteTransport) (err error) {
	//获取配送单
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:     args.ID,
		OrgID:  args.OrgID,
		InfoID: 0,
		UserID: 0,
	})
	if err != nil {
		return
	}
	//删除配送单
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "tms_transport", "id = :id AND (org_id = :org_id OR :org_id < 1)", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
	})
	if err != nil {
		return
	}
	//记录日志
	if args.Des == "" {
		args.Des = fmt.Sprint("删除配送单")
	}
	_ = appendLog(&argsAppendLog{
		OrgID:           args.OrgID,
		BindID:          args.BindID,
		TransportID:     args.ID,
		TransportBindID: data.BindID,
		Mark:            "delete",
		Des:             args.Des,
	})
	//推送MQTT
	pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 2)
	//通知取消订单
	pushNatsStatusUpdate("cancel", data.ID, "配送单被取消")
	//反馈
	return
}
