package ERPWarehouseMod

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"time"
)

// ArgsAppendLog 注入日志参数
type ArgsAppendLog struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//SN，商户下唯一，可注入其他外部系统SN
	SN string `db:"sn" json:"sn"`
	//动作类型
	// in 入库; out 出库; move_in 移动入库; move_out 移动出库
	Action string `db:"action" json:"action"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	UserID    int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//变动数量
	Count int64 `db:"count" json:"count" check:"int64Than0" empty:"true"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"600" empty:"true"`
	//附加数据，可选，如果不存在将从产品资料获取
	// 过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	// 变动时产品价格
	PerPrice int64 `db:"per_price" json:"perPrice" check:"price" empty:"true"`
}

// AppendLog 注入日志
func AppendLog(args ArgsAppendLog) {
	CoreNats.PushDataNoErr("/erp/warehouse/log_append", "", 0, "", args)
}
