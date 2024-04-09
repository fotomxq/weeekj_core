package ServiceOrderMod

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// ArgsUpdateTransportID 修改配送ID参数
type ArgsUpdateTransportID struct {
	//配送单类型
	TMSType string `json:"tmsType"`
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//编号
	// 商户下唯一
	SN int64 `db:"sn" json:"sn"`
	//今日编号
	SNDay int64 `db:"sn_day" json:"snDay"`
	//调整说明描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//配送ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id"`
}

// UpdateTransportID 修改配送ID
func UpdateTransportID(args ArgsUpdateTransportID) {
	CoreNats.PushDataNoErr("service_order_tms", "/service/order/tms", "new", args.ID, "", map[string]interface{}{
		"tmsID":   args.TransportID,
		"sn":      args.SN,
		"snDay":   args.SNDay,
		"tmsType": args.TMSType,
		"des":     args.Des,
	})
}
