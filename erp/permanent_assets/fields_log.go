package ERPPermanentAssets

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间/盘点时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//资产ID
	ProductID int64 `db:"product_id" json:"productID"`
	//模式
	// in 入库; take 领取使用; return 归还入库; check 清查库存; fix 维护; delete 销毁
	Mode string `db:"mode" json:"mode"`
	//操作主体描述
	UseName string `db:"use_name" json:"useName"`
	//实际使用人
	UseOrgBindID int64 `db:"use_org_bind_id" json:"useOrgBindID"`
	//处置后总价值
	AllPrice int64 `db:"all_price" json:"allPrice"`
	//处置后资产单价
	PerPrice int64 `db:"per_price" json:"perPrice"`
	//增加或减少数量
	Count int64 `db:"count" json:"count"`
	//存放地点
	SavePlace string `db:"save_place" json:"savePlace"`
	//备注
	Des string `db:"des" json:"des"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
