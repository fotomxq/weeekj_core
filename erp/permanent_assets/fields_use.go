package ERPPermanentAssets

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsUse struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//归还时间
	ReturnAt time.Time `db:"return_at" json:"returnAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//资产ID
	ProductID int64 `db:"product_id" json:"productID"`
	//实际使用主体（部门）
	UseName string `db:"use_name" json:"useName"`
	//实际使用人
	UseOrgBindID int64 `db:"use_org_bind_id" json:"useOrgBindID"`
	//在用数量
	Count int64 `db:"count" json:"count"`
	//领取数量
	TakeCount int64 `db:"take_count" json:"takeCount"`
	//归还数量
	ReturnCount int64 `db:"return_count" json:"returnCount"`
	//描述
	Des string `db:"des" json:"des"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
